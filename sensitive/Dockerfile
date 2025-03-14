# 多阶段构建方式

# 编译阶段
# 编译环境基础镜像
FROM golang:1.23 as stage0 
# 设置go国内代理
RUN go env -w GOPROXY=https://proxy.golang.com.cn,https://goproxy.cn,direct
# 将源码copy到临时容器
ADD ./ /src/chatgpt-sensitive
# 指定工作目录
WORKDIR /src/chatgpt-sensitive
# 交叉编译
RUN CGO_ENABLED=0 GOODS=linux GOARCH=amd64 go build -o chatgpt-sensitive .

# 运行阶段
FROM alpine:latest as stage1
WORKDIR /app/
ADD ./grpc_health_probe-linux-amd64 /app/grpc_health_probe
RUN chmod +x /app/grpc_health_probe
ADD ./dev.config.yaml /app/config.yaml
ADD ./dict.txt /app/dict.txt
COPY --from=stage0 /src/chatgpt-sensitive/chatgpt-sensitive ./

HEALTHCHECK --interval=5s --retries=3 \
CMD ./grpc_health_probe -addr=:50053

# 指定入口程序
ENTRYPOINT [ "./chatgpt-sensitive" ]
# 指定容器的启动命令或者入口程序的参数
CMD [ "--config=config.yaml","--dict=dict.txt" ]