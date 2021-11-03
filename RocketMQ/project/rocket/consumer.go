package rocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	_const "rock/const"
	"rock/dto"
)

type Consumer struct {
}

func (c Consumer) Start() {
	RocketmqConsumerClient.Subscribe(_const.GroupName, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (c consumer.ConsumeResult, err error) {
		// 执行消费的逻辑
		for i := range msgs {
			req := dto.Data{}
			err = json.Unmarshal(msgs[i].Body, &req)
			if err != nil {
				// 日志的记录
				fmt.Println(err, req)
				continue
			} else {
				//msg := string(msgs[i].Body) fmt.Println(msg)
				fmt.Println(req)
			}
		}
		c = consumer.ConsumeSuccess
		return
	})
	// 日志记录 ....
	err := RocketmqConsumerClient.Start()
	if err != nil {
		return
	}
}
