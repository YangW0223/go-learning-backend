package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/yang/go-learning-backend/examples/week11"
)

func main() {
	_ = os.Setenv("DB_DSN", "postgres://demo")
	_ = os.Setenv("JWT_SECRET", "demo-secret")
	_ = os.Setenv("APP_VERSION", "v1.1.0")

	cfg, err := week11.LoadConfigFromEnv(os.Getenv)
	if err != nil {
		fmt.Println("load config error:", err)
		return
	}
	fmt.Printf("loaded config: port=%s env=%s version=%s\n", cfg.Port, cfg.Env, cfg.Version)

	mux := week11.NewServer(cfg)
	show(mux, "/healthz")
	show(mux, "/readyz")
	show(mux, "/version")
	fmt.Println(week11.BuildRollbackPlan("demo:v2", "demo:v1"))
}

func show(mux *http.ServeMux, path string) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	fmt.Printf("GET %s => %d %s\n", path, rec.Code, rec.Body.String())
}
