# 编译阶段
# 编译环境基础镜像
FROM golang:1.23 AS stage1 
# 设置go国内代理
RUN go env -w GOPROXY=https://proxy.golang.com.cn,https://goproxy.cn,direct
# 将源码copy到临时容器
ADD ./ /src/chatgpt-wecom
# 指定工作目录
WORKDIR /src/chatgpt-wecom
# 交叉编译
RUN CGO_ENABLED=0 GOODS=linux GOARCH=amd64 go build -o chatgpt-wecom .

# 运行阶段
FROM alpine:latest AS stage2
WORKDIR /app/
ADD ./curl_x86_64_glibc_8.11.1.tar.xz /app/
HEALTHCHECK --interval=5s --retries=3 \
CMD ./curl -f http://localhost:8687/api/health || exit 1
ADD ./dev.config.yaml /app/config.yaml
COPY --from=stage1 /src/chatgpt-wecom/chatgpt-wecom /app/chatgpt-wecom


# 指定入口程序
ENTRYPOINT [ "./chatgpt-wecom" ]
# 指定容器的启动命令或者入口程序的参数
CMD [ "--config=config.yaml" ]