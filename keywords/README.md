## protoc生成代码
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/keywords.proto
```
## 构建镜像
```bash
docker build -t keywords:0.1.0 . 
```

## docker service 部署
```bash
docker config create --label env=prod keywords-conf dev.config.yaml
docker config create --label env=prod keywords-dict cainiao-coding.txt

docker service create --name keywords --config src=keywords-conf,target=/app/config.yaml --config src=keywords-dict,target=/app/dict.txt -p 50054:50054 --replicas 1 --limit-cpu 0.3 --update-parallelism=2 134.175.250.62:5000/keywords:0.1.0
```