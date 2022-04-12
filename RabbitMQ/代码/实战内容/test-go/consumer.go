package main

import (
	"log"
	"test/xmt"
)

func main() {

	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)

	rbt := &xmtmq.RabbitMQ{
		Exchange: "exchange_1",
		Key: "key_1",
		QueueName: "queue_2",
		MQUrl:     "amqp://guest:guest@127.0.0.1:5672/",
	}

	xmtmq.NewRabbitMQ(rbt)
	rbt.ConsumeRoutingMsg()
	xmtmq.RabbitMQFree(rbt)
}
