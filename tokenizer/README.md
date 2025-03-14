# tokenizer
## docker build 
```bash
docker build -t tokenizer:0.1.0 .
```
## docker 运行
```bash
# 可通过环境变量，动态指定应用程序监听端口
docker run --rm -d -p 5001:5001 -e PORT=5001 --name tokenizer tokenizer:0.1.0 
```
## docker service部署
```bash
docker tag tokenizer:0.1.0 134.175.250.62:5000/tokenizer:0.1.0
docker login -u lin  -p 123456 134.175.250.62:5000
docker push 134.175.250.62:5000/tokenizer:0.1.0

docker pull 134.175.250.62:5000/tokenizer:0.1.0
docker service create --name tokenizer --env PORT=5001 -p 5001:5001 --replicas 2 --limit-cpu 0.1 --update-parallelism=2 134.175.250.62:5000/tokenizer:0.1.0
```