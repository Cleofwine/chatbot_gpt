# wechat

# docker build
```bash
docker build -t chatgpt-wechat:0.1.0 .
```

# docker 部署服务
```bash
docker run -d -v /home/chatgpt-work/chatgpt-wechat/dev.config.yaml:/app/config.yaml --name chatgpt-wechat01 chatgpt-wechat:0.1.0
# 查看容器日志输出，获取首次扫码登录链接
docker logs chatgpt-wechat01
```