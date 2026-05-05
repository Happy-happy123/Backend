# 使用官方Go镜像作为构建阶段
FROM golang:1.24.5-alpine AS builder

# 启用国内代理，解决下载失败
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 把Go代码编译成可执行文件（main）
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/student

# 使用轻量级Alpine镜像作为运行阶段
FROM alpine:latest

# 安装ca-certificates以支持HTTPS请求
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]