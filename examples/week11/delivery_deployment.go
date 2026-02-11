// 详细注释: package week11
package week11

// 详细注释: import (
import (
	// 详细注释: "errors"
	"errors"
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "strings"
	"strings"
	// 详细注释: )
)

// 详细注释: var (
var (
	// ErrMissingConfig 表示关键配置缺失。
	// 详细注释: ErrMissingConfig = errors.New("missing required config")
	ErrMissingConfig = errors.New("missing required config")

// 详细注释: )
)

// Config 表示部署时由环境变量注入的配置。
// 详细注释: type Config struct {
type Config struct {
	// 详细注释: Port      string
	Port string
	// 详细注释: DBDSN     string
	DBDSN string
	// 详细注释: JWTSecret string
	JWTSecret string
	// 详细注释: Env       string
	Env string
	// 详细注释: Version   string
	Version string
	// 详细注释: }
}

// LoadConfigFromEnv 从 getenv 读取配置。
// 详细注释: func LoadConfigFromEnv(getenv func(string) string) (Config, error) {
func LoadConfigFromEnv(getenv func(string) string) (Config, error) {
	// 详细注释: cfg := Config{
	cfg := Config{
		// 详细注释: Port:      fallback(getenv("APP_PORT"), "8080"),
		Port: fallback(getenv("APP_PORT"), "8080"),
		// 详细注释: DBDSN:     strings.TrimSpace(getenv("DB_DSN")),
		DBDSN: strings.TrimSpace(getenv("DB_DSN")),
		// 详细注释: JWTSecret: strings.TrimSpace(getenv("JWT_SECRET")),
		JWTSecret: strings.TrimSpace(getenv("JWT_SECRET")),
		// 详细注释: Env:       fallback(getenv("APP_ENV"), "dev"),
		Env: fallback(getenv("APP_ENV"), "dev"),
		// 详细注释: Version:   fallback(getenv("APP_VERSION"), "unknown"),
		Version: fallback(getenv("APP_VERSION"), "unknown"),
		// 详细注释: }
	}
	// 详细注释: if cfg.DBDSN == "" || cfg.JWTSecret == "" {
	if cfg.DBDSN == "" || cfg.JWTSecret == "" {
		// 详细注释: return Config{}, ErrMissingConfig
		return Config{}, ErrMissingConfig
		// 详细注释: }
	}
	// 详细注释: return cfg, nil
	return cfg, nil
	// 详细注释: }
}

// NewServer 构建最小可交付服务。
// /healthz 反映进程存活，/readyz 反映依赖配置是否就绪。
// 详细注释: func NewServer(cfg Config) *http.ServeMux {
func NewServer(cfg Config) *http.ServeMux {
	// 详细注释: mux := http.NewServeMux()
	mux := http.NewServeMux()
	// 详细注释: mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		// 详细注释: w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusOK)
		// 详细注释: _, _ = w.Write([]byte("ok"))
		_, _ = w.Write([]byte("ok"))
		// 详细注释: })
	})
	// 详细注释: mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		// 详细注释: if cfg.DBDSN == "" || cfg.JWTSecret == "" {
		if cfg.DBDSN == "" || cfg.JWTSecret == "" {
			// 详细注释: http.Error(w, "not ready", http.StatusServiceUnavailable)
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusOK)
		// 详细注释: _, _ = w.Write([]byte("ready"))
		_, _ = w.Write([]byte("ready"))
		// 详细注释: })
	})
	// 详细注释: mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		// 详细注释: w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusOK)
		// 详细注释: _, _ = w.Write([]byte(cfg.Version))
		_, _ = w.Write([]byte(cfg.Version))
		// 详细注释: })
	})
	// 详细注释: return mux
	return mux
	// 详细注释: }
}

// BuildRollbackPlan 生成最小回滚说明。
// 详细注释: func BuildRollbackPlan(currentImage, previousImage string) string {
func BuildRollbackPlan(currentImage, previousImage string) string {
	// 详细注释: return fmt.Sprintf("rollback from %s to %s", currentImage, previousImage)
	return fmt.Sprintf("rollback from %s to %s", currentImage, previousImage)
	// 详细注释: }
}

// 详细注释: func fallback(value, def string) string {
func fallback(value, def string) string {
	// 详细注释: value = strings.TrimSpace(value)
	value = strings.TrimSpace(value)
	// 详细注释: if value == "" {
	if value == "" {
		// 详细注释: return def
		return def
		// 详细注释: }
	}
	// 详细注释: return value
	return value
	// 详细注释: }
}
