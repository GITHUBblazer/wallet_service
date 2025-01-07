# 使用官方的Go镜像作为基础镜像
FROM golang:1.18-alpine

# 设置工作目录
WORKDIR /app

# 复制项目的go.mod和go.sum文件到工作目录，并下载依赖包
COPY go.mod go.sum./
RUN go mod download

# 复制项目的所有源代码到工作目录
COPY..

# 构建项目可执行文件，名为main
RUN go build -o main.

# 暴露应用程序监听的端口（假设为8080）
EXPOSE 8080

# 定义容器启动时执行的命令
CMD ["./main"]