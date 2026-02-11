package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

// main 提供一个最小 migration 执行器：
// 1) 通过 -direction 选择 up/down；
// 2) 读取 migrations 目录下 SQL 文件；
// 3) 直接执行 SQL。
func main() {
	direction := flag.String("direction", "up", "migration direction: up or down")
	flag.Parse()

	// 允许通过 PG_DSN 覆盖连接串，未设置时使用本地默认值。
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/gin_backend?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("open db: %w", err))
	}
	defer db.Close()

	// 根据方向选择对应 migration 文件。
	var filename string
	switch *direction {
	case "up":
		filename = "0001_init.up.sql"
	case "down":
		filename = "0001_init.down.sql"
	default:
		panic(fmt.Errorf("unsupported direction: %s", *direction))
	}

	// 读取 SQL 文件内容并执行。
	sqlPath := filepath.Join("migrations", filename)
	queryBytes, err := os.ReadFile(sqlPath)
	if err != nil {
		panic(fmt.Errorf("read migration file: %w", err))
	}

	if _, err := db.Exec(string(queryBytes)); err != nil {
		panic(fmt.Errorf("exec migration: %w", err))
	}

	// 输出执行结果，便于在 CI 或本地终端确认。
	fmt.Printf("migration %s applied: %s\n", *direction, sqlPath)
}
