# crontab 
1. 定时刷新公众号和企业微信接口调用凭据（access_token）
2. 提供接口用于访问公众号和企微的access_token

# grpc
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/token.proto
```

# docker build
```bash
docker build -t chatgpt-crontab:0.1.0 .
```

# docker swarm 部署服务
```bash
docker config create --label env=prod chatgpt-crontab-conf dev.config.yaml

docker service create --name chatgpt-crontab --config src=chatgpt-crontab-conf,target=/app/config.yaml -p 50056:50056 --replicas 1 --limit-cpu 0.3 --update-parallelism=2 134.175.250.62:5000/chatgpt-crontab:0.1.1
```