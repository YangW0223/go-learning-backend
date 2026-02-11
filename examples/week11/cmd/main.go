// 详细注释: package main
package main

// 详细注释: import (
import (
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "net/http/httptest"
	"net/http/httptest"
	// 详细注释: "os"
	"os"

	// 详细注释: "github.com/yang/go-learning-backend/examples/week11"
	"github.com/yang/go-learning-backend/examples/week11"
	// 详细注释: )
)

// 详细注释: func main() {
func main() {
	// 详细注释: _ = os.Setenv("DB_DSN", "postgres://demo")
	_ = os.Setenv("DB_DSN", "postgres://demo")
	// 详细注释: _ = os.Setenv("JWT_SECRET", "demo-secret")
	_ = os.Setenv("JWT_SECRET", "demo-secret")
	// 详细注释: _ = os.Setenv("APP_VERSION", "v1.1.0")
	_ = os.Setenv("APP_VERSION", "v1.1.0")

	// 详细注释: cfg, err := week11.LoadConfigFromEnv(os.Getenv)
	cfg, err := week11.LoadConfigFromEnv(os.Getenv)
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: fmt.Println("load config error:", err)
		fmt.Println("load config error:", err)
		// 详细注释: return
		return
		// 详细注释: }
	}
	// 详细注释: fmt.Printf("loaded config: port=%s env=%s version=%s\n", cfg.Port, cfg.Env, cfg.Version)
	fmt.Printf("loaded config: port=%s env=%s version=%s\n", cfg.Port, cfg.Env, cfg.Version)

	// 详细注释: mux := week11.NewServer(cfg)
	mux := week11.NewServer(cfg)
	// 详细注释: show(mux, "/healthz")
	show(mux, "/healthz")
	// 详细注释: show(mux, "/readyz")
	show(mux, "/readyz")
	// 详细注释: show(mux, "/version")
	show(mux, "/version")
	// 详细注释: fmt.Println(week11.BuildRollbackPlan("demo:v2", "demo:v1"))
	fmt.Println(week11.BuildRollbackPlan("demo:v2", "demo:v1"))
	// 详细注释: }
}

// 详细注释: func show(mux *http.ServeMux, path string) {
func show(mux *http.ServeMux, path string) {
	// 详细注释: req := httptest.NewRequest(http.MethodGet, path, nil)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	// 详细注释: rec := httptest.NewRecorder()
	rec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(rec, req)
	mux.ServeHTTP(rec, req)
	// 详细注释: fmt.Printf("GET %s => %d %s\n", path, rec.Code, rec.Body.String())
	fmt.Printf("GET %s => %d %s\n", path, rec.Code, rec.Body.String())
	// 详细注释: }
}
