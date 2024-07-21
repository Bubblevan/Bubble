# 使用Go 1.13基础镜像
ARG GO_VER=1.13.8
ARG ALPINE_VER=3.10

FROM golang:${GO_VER}-alpine${ALPINE_VER}
ENV GOPROXY=https://goproxy.cn,direct

# 安装必要的依赖
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /chaincode

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制当前目录内容到工作目录
COPY . .

# 下载依赖并编译链码
RUN go mod tidy && go build -o chaincode

# 设置启动命令
EXPOSE 9999
CMD ["./chaincode"]
