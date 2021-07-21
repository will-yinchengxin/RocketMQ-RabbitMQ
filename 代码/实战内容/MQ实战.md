# MQ实战
## 1.docker 安装MQ
先拉取RabbitMQ镜像：
- docker pull rabbitmq:3.8.1-management 

然后查看镜像：
- docker images

然后运行容器(暴漏5672透出端口,及15672图像化界面端口):
- docker run --name rabbitmq -d -p 5672:5672 -p 15672:15672 -v /data:/var/lib/rabbitmq rabbitmq:3.8.1-management 

可以通过`http://localhost:15672` 直接进行访问,初始化密码/用户 均为guest

## 2.[Go](https://github.com/streadway/amqp) / [PHP](https://rabbitmq.shujuwajue.com/ying-yong-jiao-cheng/php-ban/2-work_queues.md) 使用 MQ
关于MQ发送消息的[七种模式](https://www.rabbitmq.com/tutorials/tutorial-four-php.html)

这里我们统一使用routing模式进行演示:

### 2.1 php
````
- src                 
    - RabbitMQ.php
- vender
- composer.json
- composer.lock
- index.php
````
composer.json
````
{
    "name": "willyin/amqp",
    "description": "this is the test for amqp",
    "type": "library",
    "autoload": {
        "psr-4": {
            "Willyin\\Amqp\\": "src/"
        }
    },
    "authors": [
        {
            "name": "will",
            "email": "826895143@qq.com"
        }
    ],
    "require": {
        "php-amqplib/php-amqplib": " ^2.11.0"
    }
}
````
index.php(入口文件)
````
<?php
require_once __DIR__ . '/vendor/autoload.php';

use Willyin\Amqp\RabbitMQ;

(new RabbitMQ)->push();
(new RabbitMQ)->cus();
````
RabbitMQ.php(这里采用了安装php的 [amqp](http://pecl.php.net/package/amqp) 的拓展)
````
<?php
/**
 * AMQP
 *
 * @package   Willyin\Amqp
 * @author    Will  <826895143@qq.com>
 * @copyright Copyright (C) 2021 Will
 */

namespace Willyin\Amqp;


class RabbitMQ
{
    protected $connect;

    protected $config = [
        'host' => '127.0.0.1',
        'vhost' => '/',
        'port' => 5672,
        'login' => 'guest',
        'password' => 'guest'
    ];

    /**
     * 发送消息(路由模式模式)
     *
     * 每个消费者监听自己的队列，并且设置带统配符的 routingkey
     * 生产者将消息发给broker，由交换机根据 routingkey 来转发消息到指定的队列
     */
    public function push()
    {
        $cnn = new \AMQPConnection($this->config);
        if (!$cnn->connect()) {
            echo "Can't connect to the test";
            exit();
        }
        $ch = new \AMQPChannel($cnn);
        $ex = new \AMQPExchange($ch);

        //消息的路由键，一定要和消费者端一致
        $routingKeyOne = 'key_1';
        $routingKeyTwo = 'key_2';

        //交换机名称，一定要和消费者端一致，
        $exchangeName = 'exchange_1';

        $ex->setName($exchangeName);
        $ex->setType(AMQP_EX_TYPE_DIRECT);
        $ex->setFlags(AMQP_DURABLE);
        $ex->declareExchange();

        //创建一个消息队列
        $q = new \AMQPQueue($ch);
        //设置队列名称
        $q->setName('queue_1');
        //设置队列持久
        $q->setFlags(AMQP_DURABLE);
        //声明消息队列
        $q->declareQueue();
        //交换机和队列通过$routingKey进行绑定
        $q->bind($ex->getName(), $routingKeyOne);

        /** 创建了多个queue,共享交换机中的信息
        $qOne = new \AMQPQueue($ch);
        //设置队列名称
        $qOne->setName('queue_2');
        //设置队列持久
        $qOne->setFlags(AMQP_DURABLE);
        //声明消息队列
        $qOne->declareQueue();
        //交换机和队列通过$routingKey进行绑定
        $qOne->bind($ex->getName(), $routingKeyTwo);
        **/

        //创建10个消息
        for ($i = 1; $i <= 10; $i++) {
            //消息内容
            $msg = array(
                'data' => 'message_' . $i,
                'hello' => 'world',
            );
            //发送消息到交换机，并返回发送结果
            //delivery_mode:2声明消息持久，持久的队列+持久的消息在RabbitMQ重启后才不会丢失
            echo "Send Message:" . $ex->publish(json_encode($msg), $routingKeyOne
                    , AMQP_NOPARAM, array('delivery_mode' => 2)) . "\n";
            //代码执行完毕后进程会自动退出
        }
        /**
         * 声明队列并声明交换机 -> 创建连接 -> 创建通道 -> 通道声明交换机 -> 通道声明队列 -> 通过通道使队列绑定到交换机并指定该队列的routingkey（通配符）
         *  -> 制定消息 -> 发送消息并指定routingkey（通配符）
         * */
    }

    /**
     * 消费消息
     */
    public function cus()
    {
        //连接
        $cnn = new \AMQPConnection($this->config);
        if (!$cnn->connect()) {
            echo "Cannot connect to the broker";
            exit();
        }

        //在连接内创建一个通道
        $ch = new \AMQPChannel($cnn);

        //创建一个交换机
        $ex = new \AMQPExchange($ch);

        //声明路由键
        $routingKey = 'key_1';

        //声明交换机名称
        $exchangeName = 'exchange_1';

        //设置交换机名称
        $ex->setName($exchangeName);

        //设置交换机类型
        //AMQP_EX_TYPE_DIRECT:直连交换机
        //AMQP_EX_TYPE_FANOUT:扇形交换机
        //AMQP_EX_TYPE_HEADERS:头交换机
        //AMQP_EX_TYPE_TOPIC:主题交换机
        $ex->setType(AMQP_EX_TYPE_DIRECT);

        //设置交换机持
        $ex->setFlags(AMQP_DURABLE);

        //声明交换机
        $ex->declareExchange();

        //创建一个消息队列
        $q = new \AMQPQueue($ch);

        //设置队列名称
        $q->setName('queue_1');

        //设置队列持久
        $q->setFlags(AMQP_DURABLE);

        //声明消息队列
        $q->declareQueue();

        //交换机和队列通过$routingKey进行绑定
        $q->bind($ex->getName(), $routingKey);

        //设置消息队列消费者回调方法，
        $q->consume(function ($envelope, $queue) {
            //休眠两秒，
            sleep(1);
            //echo消息内容
            echo $envelope->getBody() . "\n";
            //显式确认，队列收到消费者显式确认后，会删除该消息
            $queue->ack($envelope->getDeliveryTag());
        });
    }
}
````
### 2.2 Go
````
test
    - xmt
        - xmtmq.go
    - main.go
    - go.mod
````
main.go
````
package main

import (
	"log"
	xmtmq "test/xmt"
)

func main() {

	// 生产者模式
	//log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)
	//
	//rbt1 := &xmtmq.RabbitMQ{
	//	Exchange: "xmtPubEx2",
	//	Key: "xmt1",
	//	QueueName: "Routingqueuexmt1",
	//	MQUrl:     "amqp://guest:guest@127.0.0.1:5672/",
	//}
	//
	//xmtmq.NewRabbitMQ(rbt1)
	//rbt1.Init()
	//
	//
	//rbt2 := &xmtmq.RabbitMQ{
	//	Exchange: "xmtPubEx2",
	//	Key: "xmt2",
	//	QueueName: "Routingqueuexmt2",
	//	MQUrl:     "amqp://guest:guest@127.0.0.1:5672/xmtmq",
	//}
	//
	//xmtmq.NewRabbitMQ(rbt2)
	//rbt2.Init()
	//
	//
	//var index = 0
	//
	//for {
	//	rbt1.ProduceRouting([]byte(fmt.Sprintf("hello wolrd xmt1  %d ", index)))
	//	log.Println("发送成功xmt1  ", index)
	//
	//	rbt2.ProduceRouting([]byte(fmt.Sprintf("hello wolrd xmt2  %d ", index)))
	//	log.Println("发送成功xmt2  ", index)
	//
	//
	//	index++
	//	time.Sleep(1 * time.Second)
	//}
	//
	//
	//xmtmq.RabbitMQFree(rbt1)
	//xmtmq.RabbitMQFree(rbt2)

	// 消费者模式
	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)

	rbt := &xmtmq.RabbitMQ{
		Exchange: "exchange_1",
		Key: "key_1",
		QueueName: "queue_1",
		MQUrl:     "amqp://guest:guest@127.0.0.1:5672/",
	}

	xmtmq.NewRabbitMQ(rbt)
	rbt.ConsumeRoutingMsg()
	xmtmq.RabbitMQFree(rbt)
}
````
xmtmq.go
````
package xmtmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// routing 模式
// 定义 RabbitMQ 的数据结构
// go get github.com/streadway/amqp

type RabbitMQ struct {
	conn      *amqp.Connection // 连接
	channel   *amqp.Channel    // 通道
	QueueName string           // 队列名
	Exchange  string           // 交换机
	Key       string           // 路由键
	MQUrl     string           // MQ的虚拟机地址
}

// New 一个 RabbitMQ
func NewRabbitMQ(rbt *RabbitMQ) {
	if rbt == nil || rbt.Exchange == "" || rbt.QueueName == "" || rbt.Key == "" || rbt.MQUrl == "" {
		log.Panic("please check Exchange,QueueName,Key,MQUrl...Info")
	}

	conn, err := amqp.Dial(rbt.MQUrl)
	if err != nil {
		log.Panicf("amqp.Dial error : %v", err)
	}
	rbt.conn = conn

	channel, err := rbt.conn.Channel()
	if err != nil {
		log.Panicf("rbt.conn.Channel error : %v", err)
	}
	rbt.channel = channel
}

func RabbitMQFree(rbt *RabbitMQ) {
	if rbt == nil {
		log.Printf("rbt is nil,free failed")
		return
	}

	rbt.channel.Close()
	rbt.conn.Close()
}

func (rbt *RabbitMQ) Init() {
	// 1、创建交换机
	err := rbt.channel.ExchangeDeclare(
		rbt.Exchange, // 交换机
		amqp.ExchangeDirect,    // 交换机类型
		true,           // 是否持久化
		false,        // 是否自动删除
		false,          //  true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,          // 是否阻塞
		nil,              // 其他属性
	)
	if err != nil {
		log.Printf("rbt.channel.ExchangeDeclare error : %v", err)
		return
	}

	// 2、创建队列
	_, err = rbt.channel.QueueDeclare(
		rbt.QueueName, // 此处我们传入的是空，则是随机产生队列的名称
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("rbt.channel.QueueDeclare error : %v", err)
		return
	}

	// 3、绑定队列
	err = rbt.channel.QueueBind(
		rbt.QueueName, // 队列名字
		rbt.Key,       // routing，这里 key 需要填
		rbt.Exchange,  // 交换机名称
		false,  // 是否阻塞
		nil,      // 其他属性
	)
	if err != nil {
		log.Printf("rbt.channel.QueueBind error : %v", err)
		return
	}

}

// 生产消息 publish
func (rbt *RabbitMQ) ProduceRouting(data []byte) {
	// 1、向队列中加入数据
	err := rbt.channel.Publish(
		rbt.Exchange, 		    // 交换机
		rbt.Key,      			// key
		false,        // 若为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,        // 若为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		log.Printf("rbt.channel.Publish error : %v", err)
		return
	}

	return
}

// 消费消息
func (rbt *RabbitMQ) ConsumeRoutingMsg() {
	// 4、消费数据
	msg, err := rbt.channel.Consume(
		rbt.QueueName,    // 队列名
		"",     // 消费者的名字
		true,    // 是否自动应答
		false,  // 是否排他
		false,  // 若为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false,  // 是否阻塞
		nil,     // 其他属性
	)

	if err != nil {
		log.Printf("rbt.channel.Consume error : %v", err)
		return
	}

	for data := range msg {
		// Todo: 进行信息入库及一些消费操作
		fmt.Printf("received data is %v \n", string(data.Body))
	}

}
````
