FROM  golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

ENV GOPROXY=https://goproxy.cn,direct

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server

# 使用轻量级的alpine镜像
FROM alpine:3.21.3

# 安装ca-certificates用于HTTPS请求
RUN apk --no-cache add ca-certificates

# 创建非root用户
RUN adduser -D -h /app appuser

# 设置工作目录
WORKDIR /app

# 从builder阶段复制编译好的应用
COPY --from=builder /app/app /app/

# 切换到非root用户
USER appuser

# 暴露服务端口
EXPOSE 8080

# 运行应用
CMD ["./app"] 