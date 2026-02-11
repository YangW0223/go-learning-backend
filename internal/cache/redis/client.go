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

// Config 保存 Redis 连接参数。
type Config struct {
	Addr        string
	Password    string
	DB          int
	DialTimeout time.Duration
	IOTimeout   time.Duration
}

// Client 是一个轻量 Redis 客户端，只覆盖当前项目所需命令。
type Client struct {
	cfg Config
}

// NewClient 创建 Redis 客户端实例。
func NewClient(cfg Config) *Client {
	if strings.TrimSpace(cfg.Addr) == "" {
		cfg.Addr = "localhost:6379"
	}
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = 1 * time.Second
	}
	if cfg.IOTimeout <= 0 {
		cfg.IOTimeout = 1 * time.Second
	}

	return &Client{cfg: cfg}
}

// Ping 检查 Redis 连通性。
func (c *Client) Ping(ctx context.Context) error {
	resp, err := c.exec(ctx, "PING")
	if err != nil {
		return err
	}
	if resp.typ != '+' || resp.text != "PONG" {
		return fmt.Errorf("unexpected PING response: type=%q text=%q", string(resp.typ), resp.text)
	}
	return nil
}

// Get 读取字符串键值。
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

// SetEX 按 TTL 写入键值。
func (c *Client) SetEX(ctx context.Context, key, value string, ttl time.Duration) error {
	if ttl <= 0 {
		return fmt.Errorf("ttl must be > 0")
	}

	ttlSeconds := int(ttl / time.Second)
	if ttl%time.Second != 0 {
		ttlSeconds++
	}
	if ttlSeconds <= 0 {
		ttlSeconds = 1
	}

	resp, err := c.exec(ctx, "SETEX", key, strconv.Itoa(ttlSeconds), value)
	if err != nil {
		return err
	}
	if resp.typ != '+' || resp.text != "OK" {
		return fmt.Errorf("unexpected SETEX response: type=%q text=%q", string(resp.typ), resp.text)
	}
	return nil
}

// Del 删除指定键。
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
			return respValue{}, fmt.Errorf("unexpected AUTH response: type=%q text=%q", string(authResp.typ), authResp.text)
		}
	}

	if c.cfg.DB > 0 {
		selectResp, selectErr := c.writeAndRead(conn, "SELECT", strconv.Itoa(c.cfg.DB))
		if selectErr != nil {
			return respValue{}, selectErr
		}
		if selectResp.typ != '+' || selectResp.text != "OK" {
			return respValue{}, fmt.Errorf("unexpected SELECT response: type=%q text=%q", string(selectResp.typ), selectResp.text)
		}
	}

	return c.writeAndRead(conn, args...)
}

func (c *Client) dial(ctx context.Context) (net.Conn, error) {
	dialer := net.Dialer{Timeout: c.cfg.DialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", c.cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("dial redis %s: %w", c.cfg.Addr, err)
	}

	deadline := time.Now().Add(c.cfg.IOTimeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	if err := conn.SetDeadline(deadline); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("set redis deadline: %w", err)
	}

	return conn, nil
}

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

func writeCommand(w io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("redis command args is empty")
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
	typ     byte
	text    string
	bulk    []byte
	integer int64
	isNil   bool
	array   []respValue
}

func readRespValue(r *bufio.Reader) (respValue, error) {
	prefix, err := r.ReadByte()
	if err != nil {
		return respValue{}, fmt.Errorf("read redis response prefix: %w", err)
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
		parsed, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return respValue{}, fmt.Errorf("parse redis integer response %q: %w", line, err)
		}
		return respValue{typ: ':', integer: parsed}, nil
	case '$':
		line, err := readLine(r)
		if err != nil {
			return respValue{}, err
		}
		parsed, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return respValue{}, fmt.Errorf("parse redis bulk length %q: %w", line, err)
		}
		if parsed == -1 {
			return respValue{typ: '$', isNil: true}, nil
		}
		if parsed < 0 {
			return respValue{}, fmt.Errorf("invalid redis bulk length: %d", parsed)
		}
		payload := make([]byte, parsed+2)
		if _, err := io.ReadFull(r, payload); err != nil {
			return respValue{}, fmt.Errorf("read redis bulk payload: %w", err)
		}
		if payload[parsed] != '\r' || payload[parsed+1] != '\n' {
			return respValue{}, fmt.Errorf("invalid redis bulk payload terminator")
		}
		return respValue{typ: '$', bulk: payload[:parsed]}, nil
	case '*':
		line, err := readLine(r)
		if err != nil {
			return respValue{}, err
		}
		parsed, err := strconv.Atoi(line)
		if err != nil {
			return respValue{}, fmt.Errorf("parse redis array length %q: %w", line, err)
		}
		if parsed < 0 {
			return respValue{typ: '*', isNil: true}, nil
		}
		items := make([]respValue, 0, parsed)
		for i := 0; i < parsed; i++ {
			item, itemErr := readRespValue(r)
			if itemErr != nil {
				return respValue{}, itemErr
			}
			items = append(items, item)
		}
		return respValue{typ: '*', array: items}, nil
	default:
		return respValue{}, fmt.Errorf("unknown redis response prefix: %q", string(prefix))
	}
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read redis line: %w", err)
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	return line, nil
}
