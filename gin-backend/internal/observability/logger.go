package observability

import (
	"log"
	"os"
)

// NewStdLogger 返回标准库 logger。
// 当前输出到 stdout，便于在容器环境被日志采集系统统一收集。
func NewStdLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
}
