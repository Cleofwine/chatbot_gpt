# 纯文件系统镜像
FROM python:3.10-alpine 
# 指定工作目录
WORKDIR /app/

# 安装curl
ADD ./curl_x86_64_glibc_8.11.1.tar.xz /app/
# 使用环境变量
ENV PORT 5001


ADD tokenizer.py /app/tokenizer.py
ADD requirements.txt /app/requirements.txt
RUN pip install -i https://mirrors.aliyun.com/pypi/simple --upgrade pip
RUN pip install --root-user-action=ignore -i https://mirrors.aliyun.com/pypi/simple -r requirements.txt

# 健康检查，如果不安装curl会失败
HEALTHCHECK --interval=5s --timeout=3s --retries=5 \
CMD ./curl -f -X POST -H "Content-Type: application/json" -d '{"role":"user","content":"hello"}' http://localhost:${PORT}/tokenizer/gpt-3.5-turbo || exit 1

# 指定启动命令
CMD ["sh", "-c", "nuxt --port ${PORT} --module tokenizer.py --workers 2"]
