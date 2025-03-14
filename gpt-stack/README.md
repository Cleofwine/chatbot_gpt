```shell
docker service ls 
# 删除原先单独部署的服务，除了代理proxy
docker service rm chatgpt-crontab chatgpt-datas chatgpt-qq chatgpt-services chatgpt-webb chatgpt-wecom chatgpt-wxofficial keywords sensitive tokenizer
docker config ls 
# 删除原先单独部署的配置，除了代理proxy的
docker config rm chatgpt-crontab-conf chatgpt-data-conf chatgpt-data-conf-2 chatgpt-data-conf-3 chatgpt-proxy-conf chatgpt-qq-conf chatgpt-qq-conf-v1 chatgpt-services-conf chatgpt-webb-conf chatgpt-webb-conf-2 chatgpt-webb-conf-3 chatgpt-webb-conf-v1 chatgpt-wecom-conf chatgpt-wxofficial-conf keywords-conf keywords-dict sensitive-conf sensitive-dict

# 使用stack去部署全部服务
docker stack deploy -c compose.yaml chatgpt-stack --with-registry-auth
docker service ls 

# 如果要删除stack部署的服务，配置是不会自己删除的，需要手动删除
docker stack rm chatgpt-stack
``` 