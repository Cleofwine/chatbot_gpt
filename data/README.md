## 通过虚拟终端与容器交互
```bash
docker exec -it mysql:8.0.32 bash
```
## 命令行连接到数据
```bash
mysql -h [host] -P [port] -u root --default-character-set=utf8mb4 -p
> ./create_db.sql
> ./create_chat_records.sql
```
## protoc生成代码
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/chatgpt_data.proto
```
## 构建镜像
```bash
docker build -t chatgpt-data:0.1.0 . 
```

## docker service 部署
```bash
docker config create --label env=prod chatgpt-data-conf config.yaml

docker service create --name chatgpt-datas --config src=chatgpt-data-conf,target=/app/config.yaml -p 50052:50052 --replicas 2 --limit-cpu 0.3 --update-parallelism=2 134.175.250.62:5000/chatgpt-data:0.1.0
```