## [RocketMQ](http://rocketmq.apache.org/) 
### 一. Docker 安装
````
docker pull rocketmqinc/rocketmq // 拉取官方镜像

docker pull styletang/rocketmq-console-ng // 拉取rocketmq-console镜像(可视化界面)
````
运行(这里是windows环境)
````
// 运行nameserver
docker run -d -p 9876:9876 -v D:\\rocketMQ\\data\\namesrv\\logs:/root/logs -v D:\\rocketMQ\\data\\namesrv\\store:/root/store --name rmqnamesrv  rocketmqinc/rocketmq sh mqnamesrv
// 运行broker
docker run -d -p 10911:10911 -p 10909:10909 -v D:\\rocketMQ\\data\\broker\\logs:/root/logs -v D:\\rocketMQ\\data\\broker\\store:/root/store --name rmqbroker --link rmqnamesrv:namesrv -e "NAMESRV_ADDR=namesrv:9876" rocketmqinc/rocketmq sh mqbroker -c /opt/rocketmq-4.4.0/conf/broker.conf
````
修改配置文件: `docker exec -it rmqbroker /bin/sh`, 修改配置 `vi /opt/rocketmq-4.4.0/conf/broker.conf`
````
在broker.conf文件内新增brokerIP1属性，即：
brokerIP1 = 本机外网ip地址(ipconfig 即可获取)
````    
重启容器: `docker container restart rmqbroker`

运行rocketmq-console
````
docker run -e "JAVA_OPTS=-Drocketmq.namesrv.addr=192.168.60.66:9876 -Dcom.rocketmq.sendMessageWithVIPChannel=false" -p 8001:8080 -t styletang/rocketmq-console-ng
`````
查看主页: `http://localhost:8001`