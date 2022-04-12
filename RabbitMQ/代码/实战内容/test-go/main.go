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

