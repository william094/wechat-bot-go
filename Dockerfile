# 使用 Go 1.21 作为基础镜像
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将项目的 go.mod 和 go.sum 文件复制到容器中
COPY go.* ./ 
#COPY go.mod go.sum ./

# 下载依赖包
RUN go mod download

# 将项目的所有文件复制到工作目录
COPY . .

# 编译项目
RUN go build -o wechat-bot ./main.go

# 使用一个更小的基础镜像来运行应用程序
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制编译好的可执行文件到运行镜像中
COPY --from=builder /app/wechat-bot .

# 复制配置文件目录到容器中（如果你有默认配置文件）
COPY ./configs/wechat.yaml /app/wechat.yaml

# 设置容器启动时的入口点，并接受动态参数
ENTRYPOINT ["./wechat-bot", "run", "wechat"]

# 默认命令可以留空，由用户在 docker run 时指定
CMD ["--config-path=/app/wechat.yaml"]

# 设置容器启动时执行的命令
#CMD ["./wechat-bot", "run", "wechat", "--config-path=/app/wechat.yaml"]

