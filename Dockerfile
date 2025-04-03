# 第一阶段构建
FROM golang:1.23-alpine AS builder
# 设置构建阶段的工作目录为 /build
WORKDIR /build
# 复制整个项目目录到容器的 /build 目录

# 复制 rpc_gen 目录
COPY ./ ./

# 设置 Go 环境变量
ENV GO111MODULE=on \
    # GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux

# 构建应用
RUN go mod tidy
RUN sh build.sh

# 第二阶段运行
FROM alpine:latest
 # 设置运行阶段的工作目录为 /app
WORKDIR /app
# 从构建阶段复制文件到运行阶段的 /app 目录
COPY --from=builder /build/output/ /app/

# 设置时区
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

ENV GO_ENV=online

# grpc
EXPOSE 8080

CMD ["sh", "bootstrap.sh"]