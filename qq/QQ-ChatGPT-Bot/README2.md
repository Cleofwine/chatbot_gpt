## protoc生成代码
```bash
cd ./services
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./chatgpt-services/proto/chatgpt.proto
```
## 测试运行
```bash
go run .
```
## 构建镜像
```bash
docker image prune # 清除下虚旋镜像
docker build -t chatgpt-qq:0.1.0 . 
```

## docker service 部署
```bash
docker config create --label env=prod chatgpt-qq-conf dev.config.cfg

docker service create --name chatgpt-qq --config src=chatgpt-qq-conf,target=/app/config.cfg -p 8989:8989 --replicas 1 --limit-cpu 0.3 --update-parallelism=2 134.175.250.62:5000/chatgpt-qq:0.1.0
```