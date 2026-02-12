# 多阶段构建：
# 1) 在完整 Go 工具链镜像中编译静态二进制。
# 2) 仅将二进制拷贝到最小运行时镜像。
FROM golang:1.23-alpine AS builder

# 在构建容器内使用独立工作目录。
WORKDIR /app

# 先复制模块文件，最大化复用 Docker 构建缓存层。
COPY go.mod ./
RUN go mod download

# 复制源码并编译 Linux/amd64 静态二进制。
# 关闭 CGO，使程序可在最小 Alpine 镜像中直接运行。
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

# 运行时阶段：缩小最终镜像体积并降低攻击面。
FROM alpine:3.20

# 使用非 root 用户运行，安全默认值更好。
RUN adduser -D appuser
USER appuser
WORKDIR /app

# 仅从 builder 阶段复制编译产物。
COPY --from=builder /out/server /app/server
# 声明服务监听端口。
EXPOSE 8081

# 启动 API 服务进程。
ENTRYPOINT ["/app/server"]
