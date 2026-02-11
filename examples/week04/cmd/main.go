// package main 提供 week04 示例的独立可运行入口。
package main

// import 分组引入命令行入口所需标准库与本地 week04 包。
import (
	// flag 用于解析命令行参数。
	"flag"
	// fmt 用于打印演示输出与启动提示。
	"fmt"
	// log 用于输出服务异常并退出进程。
	"log"
	// net/http 提供 HTTP 服务与常量。
	"net/http"
	// net/http/httptest 用于在 demo 模式下本地构造请求。
	"net/http/httptest"

	// week04 包包含 handler/service/store 的分层实现。
	"github.com/yang/go-learning-backend/examples/week04"
)

// main 负责解析参数、组装依赖并选择运行模式。
func main() {
	// mode 指定运行模式：demo 本地演示；server 启动 HTTP 服务。
	mode := flag.String("mode", "demo", "run mode: demo or server")
	// addr 指定 server 模式下的监听地址。
	addr := flag.String("addr", ":18084", "listen address when mode=server")
	// Parse 触发参数解析并写回 mode/addr 指针。
	flag.Parse()

	// 创建内存 store，并预置两条待办数据。
	store := week04.NewInMemoryTodoStore([]week04.Todo{
		// 第一条用于演示正常完成流程。
		{ID: "1", Title: "read layering notes", Done: false},
		// 第二条用于展示同一套接口可操作多条数据。
		{ID: "2", Title: "write service tests", Done: false},
	})
	// 基于 store 构建 service 层。
	service := week04.NewTodoService(store)
	// 基于 service 构建 handler 层。
	handler := week04.NewTodoHandler(service)
	// 构建并获取完整路由。
	mux := week04.NewMux(handler)

	// 根据 mode 选择 server 模式或 demo 模式。
	switch *mode {
	// server 模式启动真实 HTTP 服务。
	case "server":
		runServer(*addr, mux)
	// 其他输入（含默认值 demo）走本地演示流程。
	default:
		runDemo(mux)
	}
}

// runDemo 在进程内发起三次请求，展示成功/参数错误/未找到三种结果。
func runDemo(mux *http.ServeMux) {
	// 请求合法 id=1，预期成功。
	runDemoRequest(mux, http.MethodPatch, "/api/v1/todos/1/done")
	// 请求非法 id=abc，预期参数错误。
	runDemoRequest(mux, http.MethodPatch, "/api/v1/todos/abc/done")
	// 请求不存在 id=99，预期未找到。
	runDemoRequest(mux, http.MethodPatch, "/api/v1/todos/99/done")
}

// runServer 以给定地址启动 HTTP 服务。
func runServer(addr string, mux *http.ServeMux) {
	// 输出启动地址，方便本地访问确认。
	fmt.Printf("week04 server listening on http://localhost%s\n", addr)
	// 输出可调用的接口路径说明。
	fmt.Println("endpoint: PATCH /api/v1/todos/{id}/done")
	// 启动监听；若异常退出则打印错误并终止程序。
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

// runDemoRequest 构造一次请求并打印状态码与响应体。
func runDemoRequest(mux *http.ServeMux, method, path string) {
	// 构造一个无请求体的测试请求。
	req := httptest.NewRequest(method, path, nil)
	// 创建响应记录器用于捕获 handler 输出。
	recorder := httptest.NewRecorder()
	// 执行一次完整路由处理流程。
	mux.ServeHTTP(recorder, req)

	// 打印请求方法、路径、状态码和响应体，便于学习观察。
	fmt.Printf("%s %s => status=%d body=%s", method, path, recorder.Code, recorder.Body.String())
}
