# wxofficial 
# grpc
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/token.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/chatgpt.proto
```

# docker build
```bash
docker build -t chatgpt-wxofficial:0.1.0 .
```

# docker swarm 部署服务
```bash
docker config create --label env=prod chatgpt-wxofficial-conf dev.config.yaml

docker service create --name chatgpt-wxofficial --config src=chatgpt-wxofficial-conf,target=/app/config.yaml -p 8686:8686 --replicas 1 --limit-cpu 0.3 --update-parallelism=2 134.175.250.62:5000/chatgpt-wxofficial:0.1.1
```