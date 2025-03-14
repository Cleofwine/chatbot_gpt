# 多阶段构建方式
# 编译阶段
# 编译环境基础镜像
FROM golang:1.23 AS stage0 
# 设置go国内代理
RUN go env -w GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
# 将源码copy到临时容器
ADD ./ /src/chatgpt-qq
# 指定工作目录
WORKDIR /src/chatgpt-qq
# 交叉编译
RUN CGO_ENABLED=0 GOODS=linux GOARCH=amd64 go build -o chatgpt-qq .

# 运行阶段
FROM alpine:latest AS stage1
WORKDIR /app/
ADD ./curl_x86_64_glibc_8.11.1.tar.xz /app/
HEALTHCHECK --interval=5s --retries=3 \
CMD ./curl -f http://localhost:8989/health || exit 1
ADD ./dev.config.cfg /app/config.cfg
COPY --from=stage0 /src/chatgpt-qq/chatgpt-qq /app/chatgpt-qq


# 指定入口程序
ENTRYPOINT [ "./chatgpt-qq" ]