package rocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"log"
	_const "rock/const"
	"rock/dto"
)

type Consumer struct {
}

func (c Consumer) PushConsumerStart() {

	RocketmqPushConsumerClient.Subscribe(_const.Topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (c consumer.ConsumeResult, err error) {
		// 执行消费的逻辑, 利用分布式锁来避免重复消费
		for i := range msgs {
			req := dto.Data{}
			err = json.Unmarshal(msgs[i].Body, &req)
			if err != nil {
				// 日志的记录到 db， 或者 es 等中间件
				fmt.Println(err, req)
				continue
			}
			fmt.Println("willXXXX", string(msgs[i].Body))

			// 利用分布式锁保证 幂等性
			//redis.Redis.HSET(msgs[i].MsgId, string(msgs[i].Body))
		}

		c = consumer.ConsumeSuccess
		return
	})

	err := RocketmqPushConsumerClient.Start()
	if err != nil {
		fmt.Printf("consume data error: %s", err.Error())
	}
	log.Print("PushConsumerStart success")
	return
}

/*
目前 go client api 还未支持
*/
//func (c Consumer) PullConsumerStart() {
//	ctx := context.Background()
//	queue := primitive.MessageQueue{
//		Topic:      _const.Topic,
//		BrokerName: _const.BrokerName,
//		QueueId:    0,
//	}
//
//	offset := int64(0)
//	for {
//		/*
//			参数三为 消费者消费的的偏移量，即消费位点（）
//			参数四为每次最大读取值，摩默认最大值为 32
//		*/
//		resp, err := RocketmqPullConsumerClient.PullFrom(ctx, queue, offset, _const.MaxNumbers)
//		if err != nil {
//			if err == rocketmq.ErrRequestTimeout {
//				fmt.Printf("timeout \n")
//				time.Sleep(1 * time.Second)
//				continue
//			}
//			fmt.Printf("pull consumer get message err: %v \n", err)
//			return
//		}
//		if resp.Status == primitive.PullFound {
//			fmt.Printf("pull message success. nextOffset: %d \n", resp.NextBeginOffset)
//			for _, msg := range resp.GetMessageExts() {
//				fmt.Printf("pull msg: %v \n", msg)
//			}
//		}
//		offset = resp.NextBeginOffset
//	}
//}
