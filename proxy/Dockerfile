# 多阶段构建方式

# 编译阶段
# 编译环境基础镜像
FROM golang:1.23 as stage0 
# 设置go国内代理
RUN go env -w GOPROXY=https://proxy.golang.com.cn,https://goproxy.cn,direct
# 将源码copy到临时容器
ADD ./ /src/chatgpt-proxy
# 指定工作目录
WORKDIR /src/chatgpt-proxy
# 交叉编译
RUN CGO_ENABLED=0 GOODS=linux GOARCH=amd64 go build -o chatgpt-proxy .

# 运行阶段
FROM alpine:latest as stage1
# 健康检查设置
ADD ./curl_x86_64_glibc_8.11.1.tar.xz /app/
HEALTHCHECK --interval=5s --retries=3 \
CMD ./curl -f http://localhost:4002/health || exit 1

LABEL author="lin"
WORKDIR /app/
ADD ./config.yaml.tmp /app/config.yaml
COPY --from=stage0 /src/chatgpt-proxy/chatgpt-proxy ./
# 指定入口程序
ENTRYPOINT [ "./chatgpt-proxy" ]
# 指定容器的启动命令或者入口程序的参数
CMD [ "--config=config.yaml" ]

# docker build --no-cache -t chatgpt-proxy:latest -t chatgpt-proxy:0.1.0 .
# 产生三个镜像，一个是非常大的第一个阶段的镜像，一个是0.1.0版本的，一个是latest版本的