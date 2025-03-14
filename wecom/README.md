# wecom

# docker build
```bash
docker build -t chatgpt-wecom:0.1.1 .
```

# docker swarm 部署服务
```bash
docker config create --label env=prod chatgpt-wecom-conf dev.config.yaml

docker service create --name chatgpt-wecom --config src=chatgpt-wecom-conf,target=/app/config.yaml -p 8687:8687 --replicas 1 --limit-cpu 0.3 --update-parallelism=2 134.175.250.62:5000/chatgpt-wecom:0.1.0
```