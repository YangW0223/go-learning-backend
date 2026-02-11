package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config 是应用运行所需的聚合配置。
type Config struct {
	Server ServerConfig
	Redis  RedisConfig
}

// ServerConfig 保存 HTTP 服务监听配置。
type ServerConfig struct {
	Port string
}

// RedisConfig 保存 Redis 连接与缓存相关配置。
type RedisConfig struct {
	Enabled     bool
	Addr        string
	Password    string
	DB          int
	CacheTTL    time.Duration
	DialTimeout time.Duration
	IOTimeout   time.Duration
}

const (
	defaultPort             = "8080"
	defaultRedisEnabled     = false
	defaultRedisAddr        = "localhost:6379"
	defaultRedisPassword    = ""
	defaultRedisDB          = 0
	defaultRedisCacheTTL    = 30 * time.Second
	defaultRedisDialTimeout = 1 * time.Second
	defaultRedisIOTimeout   = 1 * time.Second
)

// Load 从环境变量加载配置，并在必要时返回可读错误。
func Load() (Config, error) {
	cfg := Config{}
	cfg.Server.Port = getEnvWithDefault("PORT", defaultPort)

	redisEnabled, err := getBoolEnv("REDIS_ENABLED", defaultRedisEnabled)
	if err != nil {
		return Config{}, err
	}

	redisDB, err := getIntEnv("REDIS_DB", defaultRedisDB)
	if err != nil {
		return Config{}, err
	}

	cacheTTLSeconds, err := getIntEnv("REDIS_CACHE_TTL_SECONDS", int(defaultRedisCacheTTL.Seconds()))
	if err != nil {
		return Config{}, err
	}
	if cacheTTLSeconds <= 0 {
		return Config{}, fmt.Errorf("REDIS_CACHE_TTL_SECONDS must be > 0")
	}

	dialTimeoutMS, err := getIntEnv("REDIS_DIAL_TIMEOUT_MS", int(defaultRedisDialTimeout.Milliseconds()))
	if err != nil {
		return Config{}, err
	}
	if dialTimeoutMS <= 0 {
		return Config{}, fmt.Errorf("REDIS_DIAL_TIMEOUT_MS must be > 0")
	}

	ioTimeoutMS, err := getIntEnv("REDIS_IO_TIMEOUT_MS", int(defaultRedisIOTimeout.Milliseconds()))
	if err != nil {
		return Config{}, err
	}
	if ioTimeoutMS <= 0 {
		return Config{}, fmt.Errorf("REDIS_IO_TIMEOUT_MS must be > 0")
	}

	redisAddr := defaultRedisAddr
	redisAddrRaw, redisAddrExists := os.LookupEnv("REDIS_ADDR")
	if redisAddrExists {
		redisAddr = strings.TrimSpace(redisAddrRaw)
	}
	if redisEnabled && redisAddr == "" {
		return Config{}, fmt.Errorf("REDIS_ADDR is required when REDIS_ENABLED=true")
	}
	if redisAddr == "" {
		redisAddr = defaultRedisAddr
	}

	cfg.Redis = RedisConfig{
		Enabled:     redisEnabled,
		Addr:        redisAddr,
		Password:    getEnvWithDefault("REDIS_PASSWORD", defaultRedisPassword),
		DB:          redisDB,
		CacheTTL:    time.Duration(cacheTTLSeconds) * time.Second,
		DialTimeout: time.Duration(dialTimeoutMS) * time.Millisecond,
		IOTimeout:   time.Duration(ioTimeoutMS) * time.Millisecond,
	}

	return cfg, nil
}

func getEnvWithDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getBoolEnv(key string, fallback bool) (bool, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be bool: %w", key, err)
	}
	return parsed, nil
}

func getIntEnv(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be int: %w", key, err)
	}
	return parsed, nil
}
