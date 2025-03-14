## protoc生成代码
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/sensitive.proto
```
## 构建镜像
```bash
docker build -t sensitive:0.1.0 . 
```

## docker service 部署
```bash
docker config create --label env=prod sensitive-conf dev.config.yaml
docker config create --label env=prod sensitive-dict dict.txt

docker service create --name sensitive --config src=sensitive-conf,target=/app/config.yaml --config src=sensitive-dict,target=/app/dict.txt -p 50053:50053 --replicas 1 --limit-cpu 0.3 --update-parallelism=2 134.175.250.62:5000/sensitive:0.1.0
```