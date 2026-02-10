package handler

import "net/http"

// Ping 是最小可用的连通性接口。
// 它不依赖数据库或其他外部资源，只返回固定响应，
// 常用于快速判断服务是否正常对外提供 HTTP。
func Ping(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "pong"})
}
