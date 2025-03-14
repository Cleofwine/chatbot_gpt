## protoc生成代码
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/chatgpt.proto
cd ./services
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./keywords/proto/keywords.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./sensitive-words/proto/sensitive.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./datas/proto/chatgpt_data.proto
```
## 构建镜像
```bash
docker build -t chatgpt-services:0.1.0 . 
```

## docker service 部署
```bash
docker config create --label env=prod chatgpt-services-conf dev.config.yaml

docker service create --name chatgpt-services --config src=chatgpt-services-conf,target=/app/config.yaml -p 50051:50051 --replicas 2 --limit-cpu 0.3 --update-parallelism=2 134.175.250.62:5000/chatgpt-services:0.1.0
```