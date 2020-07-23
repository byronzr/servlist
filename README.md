# servlist
+ 基于redis的单机docker服务列表，全自动注册。
+ host机内，以非docker方式安装了redis，docker内的微服务，以172.17.0.1的网段，访问6379将docker内的微服务的172.17.*记录在redis中。
+ 键值以15秒ttl超期，servlist以12秒的方式重续期。

## 仅服务于host机内多个docker，如要集群k8s 
