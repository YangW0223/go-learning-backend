package redis

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// Config 保存 Redis 客户端参数。
type Config struct {
	Addr             string
	Password         string
	DB               int
	DialTimeout      time.Duration
	ReadWriteTimeout time.Duration
}

// Client 是轻量 Redis 客户端，仅实现当前项目所需命令。
type Client struct {
	cfg Config
}

// NewClient 创建 Redis 客户端，并处理默认配置回填。
func NewClient(cfg Config) *Client {
	if strings.TrimSpace(cfg.Addr) == "" {
		cfg.Addr = "localhost:6379"
	}
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = 1 * time.Second
	}
	if cfg.ReadWriteTimeout <= 0 {
		cfg.ReadWriteTimeout = 1 * time.Second
	}
	return &Client{cfg: cfg}
}

// Ping 执行 Redis PING 命令，验证连通性。
func (c *Client) Ping(ctx context.Context) error {
	resp, err := c.exec(ctx, "PING")
	if err != nil {
		return err
	}
	if resp.typ != '+' || resp.text != "PONG" {
		return fmt.Errorf("unexpected ping response: %q", resp.text)
	}
	return nil
}

// Get 获取字符串键。
// 返回值约定：
// 1) hit=true 代表命中；
// 2) hit=false 代表 key 不存在；
// 3) error!=nil 代表执行失败。
func (c *Client) Get(ctx context.Context, key string) (string, bool, error) {
	resp, err := c.exec(ctx, "GET", key)
	if err != nil {
		return "", false, err
	}
	if resp.typ == '$' && resp.isNil {
		return "", false, nil
	}
	if resp.typ != '$' {
		return "", false, fmt.Errorf("unexpected GET response type: %q", string(resp.typ))
	}
	return string(resp.bulk), true, nil
}

// SetEX 设置带过期时间的字符串键。
func (c *Client) SetEX(ctx context.Context, key string, value string, ttl time.Duration) error {
	seconds := int(ttl / time.Second)
	if ttl%time.Second != 0 {
		seconds++
	}
	if seconds <= 0 {
		seconds = 1
	}
	resp, err := c.exec(ctx, "SETEX", key, strconv.Itoa(seconds), value)
	if err != nil {
		return err
	}
	if resp.typ != '+' || resp.text != "OK" {
		return fmt.Errorf("unexpected SETEX response: %q", resp.text)
	}
	return nil
}

// Del 删除指定键（key 不存在也视为成功）。
func (c *Client) Del(ctx context.Context, key string) error {
	resp, err := c.exec(ctx, "DEL", key)
	if err != nil {
		return err
	}
	if resp.typ != ':' {
		return fmt.Errorf("unexpected DEL response type: %q", string(resp.typ))
	}
	return nil
}

// exec 建立短连接并执行一条命令。
// 该实现强调简单性，适用于学习项目与低流量场景。
func (c *Client) exec(ctx context.Context, args ...string) (respValue, error) {
	conn, err := c.dial(ctx)
	if err != nil {
		return respValue{}, err
	}
	defer conn.Close()

	if c.cfg.Password != "" {
		authResp, authErr := c.writeAndRead(conn, "AUTH", c.cfg.Password)
		if authErr != nil {
			return respValue{}, authErr
		}
		if authResp.typ != '+' || authResp.text != "OK" {
			return respValue{}, fmt.Errorf("unexpected AUTH response")
		}
	}
	if c.cfg.DB > 0 {
		selectResp, selectErr := c.writeAndRead(conn, "SELECT", strconv.Itoa(c.cfg.DB))
		if selectErr != nil {
			return respValue{}, selectErr
		}
		if selectResp.typ != '+' || selectResp.text != "OK" {
			return respValue{}, fmt.Errorf("unexpected SELECT response")
		}
	}
	return c.writeAndRead(conn, args...)
}

// dial 建立 TCP 连接并设置读写超时。
func (c *Client) dial(ctx context.Context) (net.Conn, error) {
	dialer := net.Dialer{Timeout: c.cfg.DialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", c.cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("dial redis: %w", err)
	}
	deadline := time.Now().Add(c.cfg.ReadWriteTimeout)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	if err := conn.SetDeadline(deadline); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("set redis deadline: %w", err)
	}
	return conn, nil
}

// writeAndRead 执行一次“写命令 + 读响应”流程。
func (c *Client) writeAndRead(conn net.Conn, args ...string) (respValue, error) {
	if err := writeCommand(conn, args...); err != nil {
		return respValue{}, err
	}
	reader := bufio.NewReader(conn)
	resp, err := readRespValue(reader)
	if err != nil {
		return respValue{}, err
	}
	return resp, nil
}

// writeCommand 把命令参数编码为 RESP 协议格式。
func writeCommand(w io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("empty redis args")
	}
	var buf bytes.Buffer
	if _, err := fmt.Fprintf(&buf, "*%d\r\n", len(args)); err != nil {
		return err
	}
	for _, arg := range args {
		if _, err := fmt.Fprintf(&buf, "$%d\r\n%s\r\n", len(arg), arg); err != nil {
			return err
		}
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("write redis command: %w", err)
	}
	return nil
}

type respValue struct {
	typ   byte
	text  string
	bulk  []byte
	isNil bool
}

// readRespValue 解析一条 RESP 响应。
func readRespValue(r *bufio.Reader) (respValue, error) {
	prefix, err := r.ReadByte()
	if err != nil {
		return respValue{}, fmt.Errorf("read redis prefix: %w", err)
	}
	switch prefix {
	case '+':
		line, err := readLine(r)
		if err != nil {
			return respValue{}, err
		}
		return respValue{typ: '+', text: line}, nil
	case '-':
		line, err := readLine(r)
		if err != nil {
			return respValue{}, err
		}
		return respValue{}, fmt.Errorf("redis error: %s", line)
	case ':':
		line, err := readLine(r)
		if err != nil {
			return respValue{}, err
		}
		return respValue{typ: ':', text: line}, nil
	case '$':
		line, err := readLine(r)
		if err != nil {
			return respValue{}, err
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return respValue{}, fmt.Errorf("parse bulk size: %w", err)
		}
		if n == -1 {
			return respValue{typ: '$', isNil: true}, nil
		}
		if n < 0 {
			return respValue{}, fmt.Errorf("invalid bulk size")
		}
		payload := make([]byte, n+2)
		if _, err := io.ReadFull(r, payload); err != nil {
			return respValue{}, fmt.Errorf("read bulk payload: %w", err)
		}
		return respValue{typ: '$', bulk: payload[:n]}, nil
	default:
		return respValue{}, fmt.Errorf("unexpected redis prefix: %q", string(prefix))
	}
}

// readLine 读取 CRLF 结尾的一行文本。
func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read line: %w", err)
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	return line, nil
}
