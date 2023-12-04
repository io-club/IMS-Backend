FROM golang:1.21.4-alpine3.18

# 镜像换源
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.18/main/ > /etc/apk/repositories
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.18/community/ >> /etc/apk/repositories

# 安装所需工具
RUN apk add --no-cache bash make

WORKDIR /IO-IMS
# 将当前目录下所有文件拷贝到工作目录下
ADD . ./

# 设置Go语言的环境变量
# 打开Go Module模式
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn \
    mode=release

# 下载依赖
RUN go mod download

# 编译
RUN make build

# 暴露端口
EXPOSE 10000

# 启动服务，并提供终端
ENTRYPOINT ["/bin/sh", "-c", "make run && sh"]


