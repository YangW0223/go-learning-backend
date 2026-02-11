package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config 聚合应用全部配置。
type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Log      LogConfig
}

// AppConfig 保存应用级配置。
type AppConfig struct {
	Name string
	Env  string
}

// HTTPConfig 保存 HTTP 监听和超时配置。
type HTTPConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
	TrustedProxies  []string
}

// PostgresConfig 保存 Postgres 连接配置。
type PostgresConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig 保存 Redis 连接与缓存配置。
type RedisConfig struct {
	Enabled          bool
	Addr             string
	Password         string
	DB               int
	DialTimeout      time.Duration
	ReadWriteTimeout time.Duration
	CacheTTL         time.Duration
}

// AuthConfig 保存鉴权相关配置。
type AuthConfig struct {
	JWTSecret       string
	TokenTTL        time.Duration
	Issuer          string
	PasswordMinSize int
}

// LogConfig 保存日志配置。
type LogConfig struct {
	Level string
}

const (
	defaultAppName            = "gin-backend"
	defaultAppEnv             = "dev"
	defaultHTTPPort           = "8081"
	defaultHTTPReadTimeoutMS  = 5000
	defaultHTTPWriteTimeoutMS = 10000
	defaultShutdownTimeoutMS  = 10000
	defaultRequestTimeoutMS   = 3000
	defaultPGMaxOpenConns     = 20
	defaultPGMaxIdleConns     = 10
	defaultPGConnLifeMS       = 300000
	defaultRedisAddr          = "localhost:6379"
	defaultRedisDB            = 0
	defaultRedisDialMS        = 1000
	defaultRedisRWMS          = 1000
	defaultRedisCacheTTL      = 30
	defaultTokenTTLMinutes    = 120
	defaultAuthIssuer         = "gin-backend"
	defaultLogLevel           = "info"
	defaultPasswordMinSize    = 8
)

// Load 从环境变量加载配置并做必要校验。
func Load() (Config, error) {
	// cfg 是最终返回的配置对象，后续按模块逐段填充。
	cfg := Config{}

	// 加载应用级配置（服务名、运行环境）。
	cfg.App = AppConfig{
		Name: getenv("APP_NAME", defaultAppName),
		Env:  getenv("APP_ENV", defaultAppEnv),
	}

	// 加载 HTTP 超时配置（单位毫秒），随后统一转换成 time.Duration。
	httpReadMS, err := getInt("HTTP_READ_TIMEOUT_MS", defaultHTTPReadTimeoutMS)
	if err != nil {
		return Config{}, err
	}
	httpWriteMS, err := getInt("HTTP_WRITE_TIMEOUT_MS", defaultHTTPWriteTimeoutMS)
	if err != nil {
		return Config{}, err
	}
	shutdownMS, err := getInt("HTTP_SHUTDOWN_TIMEOUT_MS", defaultShutdownTimeoutMS)
	if err != nil {
		return Config{}, err
	}
	requestMS, err := getInt("HTTP_REQUEST_TIMEOUT_MS", defaultRequestTimeoutMS)
	if err != nil {
		return Config{}, err
	}

	// 组装 HTTP 配置，TrustedProxies 支持逗号分隔。
	cfg.HTTP = HTTPConfig{
		Port:            getenv("HTTP_PORT", defaultHTTPPort),
		ReadTimeout:     time.Duration(httpReadMS) * time.Millisecond,
		WriteTimeout:    time.Duration(httpWriteMS) * time.Millisecond,
		ShutdownTimeout: time.Duration(shutdownMS) * time.Millisecond,
		RequestTimeout:  time.Duration(requestMS) * time.Millisecond,
		TrustedProxies:  splitCSV(getenv("HTTP_TRUSTED_PROXIES", "")),
	}

	// 读取 Postgres 连接池参数。
	maxOpen, err := getInt("PG_MAX_OPEN_CONNS", defaultPGMaxOpenConns)
	if err != nil {
		return Config{}, err
	}
	maxIdle, err := getInt("PG_MAX_IDLE_CONNS", defaultPGMaxIdleConns)
	if err != nil {
		return Config{}, err
	}
	connLifeMS, err := getInt("PG_CONN_MAX_LIFETIME_MS", defaultPGConnLifeMS)
	if err != nil {
		return Config{}, err
	}

	// 组装 Postgres 配置。
	cfg.Postgres = PostgresConfig{
		DSN:             getenv("PG_DSN", "postgres://postgres:postgres@localhost:5432/gin_backend?sslmode=disable"),
		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxLifetime: time.Duration(connLifeMS) * time.Millisecond,
	}

	// 读取 Redis 开关与连接参数。
	redisEnabled, err := getBool("REDIS_ENABLED", true)
	if err != nil {
		return Config{}, err
	}
	redisDB, err := getInt("REDIS_DB", defaultRedisDB)
	if err != nil {
		return Config{}, err
	}
	redisDialMS, err := getInt("REDIS_DIAL_TIMEOUT_MS", defaultRedisDialMS)
	if err != nil {
		return Config{}, err
	}
	redisRWMS, err := getInt("REDIS_RW_TIMEOUT_MS", defaultRedisRWMS)
	if err != nil {
		return Config{}, err
	}
	redisCacheTTLSeconds, err := getInt("REDIS_CACHE_TTL_SECONDS", defaultRedisCacheTTL)
	if err != nil {
		return Config{}, err
	}

	// 组装 Redis 配置（包含缓存 TTL）。
	cfg.Redis = RedisConfig{
		Enabled:          redisEnabled,
		Addr:             getenv("REDIS_ADDR", defaultRedisAddr),
		Password:         getenv("REDIS_PASSWORD", ""),
		DB:               redisDB,
		DialTimeout:      time.Duration(redisDialMS) * time.Millisecond,
		ReadWriteTimeout: time.Duration(redisRWMS) * time.Millisecond,
		CacheTTL:         time.Duration(redisCacheTTLSeconds) * time.Second,
	}

	// 读取鉴权配置（JWT 过期时间与密码最小长度）。
	tokenTTLMin, err := getInt("TOKEN_TTL_MINUTES", defaultTokenTTLMinutes)
	if err != nil {
		return Config{}, err
	}
	passwordMin, err := getInt("PASSWORD_MIN_LENGTH", defaultPasswordMinSize)
	if err != nil {
		return Config{}, err
	}

	// 组装鉴权配置。
	cfg.Auth = AuthConfig{
		JWTSecret:       getenv("JWT_SECRET", "change-me"),
		TokenTTL:        time.Duration(tokenTTLMin) * time.Minute,
		Issuer:          getenv("JWT_ISSUER", defaultAuthIssuer),
		PasswordMinSize: passwordMin,
	}

	// 日志级别统一转成小写，便于后续分支判断。
	cfg.Log = LogConfig{Level: strings.ToLower(getenv("LOG_LEVEL", defaultLogLevel))}

	// 返回前做兜底校验，启动阶段尽早失败。
	if err := validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// validate 执行关键字段校验，避免服务在运行时才暴露配置错误。
func validate(cfg Config) error {
	if strings.TrimSpace(cfg.HTTP.Port) == "" {
		return fmt.Errorf("HTTP_PORT is required")
	}
	if strings.TrimSpace(cfg.Postgres.DSN) == "" {
		return fmt.Errorf("PG_DSN is required")
	}
	if cfg.Redis.Enabled && strings.TrimSpace(cfg.Redis.Addr) == "" {
		return fmt.Errorf("REDIS_ADDR is required when REDIS_ENABLED=true")
	}
	if len(cfg.Auth.JWTSecret) < 8 {
		return fmt.Errorf("JWT_SECRET length must be >= 8")
	}
	if cfg.Auth.PasswordMinSize < 6 {
		return fmt.Errorf("PASSWORD_MIN_LENGTH must be >= 6")
	}
	return nil
}

// getenv 读取字符串环境变量；当值为空时返回默认值。
func getenv(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

// getInt 读取整型环境变量；为空时使用默认值。
func getInt(key string, fallback int) (int, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("%s must be int: %w", key, err)
	}
	return parsed, nil
}

// getBool 读取布尔环境变量；为空时使用默认值。
func getBool(key string, fallback bool) (bool, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("%s must be bool: %w", key, err)
	}
	return parsed, nil
}

// splitCSV 把逗号分隔字符串转换成字符串切片，并过滤空值。
func splitCSV(v string) []string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}
