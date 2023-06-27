# 基于 golang 镜像构建
FROM golang:1.20-alpine AS build
WORKDIR /app

# 复制代码到容器内
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -o exporter .

# 使用 alpine 镜像作为最终的镜像
FROM alpine:latest

# 拷贝二进制文件
COPY --from=build /app/exporter /exporter

# 设置容器入口命令
CMD ["/exporter"]