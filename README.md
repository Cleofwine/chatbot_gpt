# 简介
- Chatgpt多端接入平台。

# 服务结构
![服务结构](https://bbucket-1253575676.cos.ap-guangzhou.myqcloud.com/public/%E6%97%A0%E6%A0%87%E9%A2%98-2025-03-13-2301.png "结构图")

# 部署
1. 前置准备
```bash
# 1. mysql
mysql -u 用户名 -p 
source ./data/sql/create_db.sql;
USE chatGPT;
source ./data/sql/create_chat_records.sql;
# 2. 确保redis启用
```
2. 构建镜像
```bash
cd ./crontab
docker build -t chatgpt-crontab:0.1.0 .
cd ./data
docker build -t chatgpt-datas:0.1.0 .
cd ./keywords
docker build -t keywords:0.1.0 .
cd ./proxy
docker build -t chatgpt-proxy:0.1.0 .
cd ./qq/QQ-ChatGPT-Bot
docker build -t chatgpt-qq:0.1.0 .
cd ./sensitive
docker build -t sensitive:0.1.0 .
cd ./service
docker build -t chatgpt-services:0.1.0 .
cd ./tokenizer
docker build -t tokenizer:0.1.0 .
cd ./web-backend
docker build -t chatgpt-webb:0.1.0 .
cd ./web-frontend
docker build -t chatgpt-frontend:0.1.0 .
cd ./wecom
docker build -t chatgpt-wecom:0.1.0 .
cd ./wxofficial
docker build -t chatgpt-wxofficial:0.1.0 .
cd ./wechat
docker build -t chatgpt-wechat:0.1.0 .
```
3. 修改配置
```bash
# 集群配置
cd ./stack/configs # 这个路径下可以找到全部的配置，修改成自己的对应COS token、redis、mysql连接方式
# 独立服务配置
cd ./proxy/ # 独立部署的代理配置在这里
cd ./qq/go-cqhttp # qq账号配置位置
``` 
4. 部署独立服务
```bash
# 境外代理
cd ./proxy # 参考READMD.md部署境外代理
# 部署qq的cqhttp
cd ./qq/go-cqhttp 
./go-cqhttp
# 部署微信bot
cd ./wechat # 参考READMD.md部署，不采用统一部署是因为单个微信号容易被封
```
5. 统一部署服务
```bash
cd ./gpt-stack
docker stack deploy -c compose.yaml aichatbot-stack
# 访问
http://<IP>:7070/webb/
```
5. 统一删除服务
```bash
# 配置不会被删除
docker stack rm aichatbot-stack
```

# 参考项目
```bash
对各个开源项目的作者们表示诚挚的感谢。
```
- https://github.com/importcjj/sensitive
- https://github.com/eatmoreapple/openwechat
- https://github.com/LagrangeDev/go-cqhttp
- https://github.com/SuInk/QQ-ChatGPT-Bot
- https://github.com/Chanzhaoyu/chatgpt-web