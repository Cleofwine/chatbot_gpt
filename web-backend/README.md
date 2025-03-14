## protoc生成代码
```bash
cd ./services
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./chatgpt-services/proto/chatgpt.proto
```
## 测试运行
```bash
go run ./cmd/main.go --config=dev.config.yaml
```
## 构建镜像
```bash
docker image prune # 清除下虚旋镜像
# 1. 需要先构建前端docker镜像
# docker build -t chatgpt-web-frontend:0.1.1 . 
# 2. 构建后端
# docker build -t chatgpt-webb:0.1.4 --build-arg "frontend_img=chatgpt-web-frontend:0.1.1" .
docker build -t chatgpt-webb:0.2.0 .
```

## docker service 部署
```bash
docker config create --label env=prod chatgpt-webb-conf-6 dev.config.yaml

docker service create --name chatgpt-webb --config src=chatgpt-webb-conf-3,target=/app/config.yaml -p 7080:7080 --replicas 2 --limit-cpu 0.3 --update-parallelism=2 134.175.250.62:5000/chatgpt-webb:0.1.0
```