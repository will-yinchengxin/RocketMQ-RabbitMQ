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

