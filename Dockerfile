#源镜像
FROM golang:latest
#作者
MAINTAINER lemonshwang "lemonshwang@tencent.com"
#设置工作目录
WORKDIR /root/learn
#将服务器的go工程代码加入到docker容器中
ADD . /root/learn
#go构建可执行文件
#RUN go build -o server server.go
RUN chmod u+x server
#暴露端口
EXPOSE 8080
#最终运行docker的命令
ENTRYPOINT  ["./server"]