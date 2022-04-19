package rocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"os"
	"reflect"
	_const "rock/const"
	"rock/dto"
)

type Producer struct {
}

func (p *Producer) Send() (res interface{}, err error) {
	data := dto.GetData()
	v := reflect.ValueOf(data)
	var sendData []*primitive.Message
	for i := 0; i < len(data); i++ {
		value := v.Index(i) // Value of item
		infoByte, _ := json.Marshal(value.Interface())

		msg := &primitive.Message{
			Topic: _const.Topic,
			Body:  infoByte,
			Queue: &primitive.MessageQueue{
				Topic:      _const.Topic,
				BrokerName: _const.BrokerName,
				QueueId:    2,
			},
		}
		msg.WithDelayTimeLevel(3)
		sendData = append(sendData, msg)
	}

	/*
		发送方式：
			- 同步发送 SendSync
			- 异步发送 SendAsync
			- 单向发送 SendOneWay

		这里采用的是异步发送的方式
	*/
	ctx := context.TODO() // 这里可以做成外部传递的形式
	err = RocketmqProducerClient.SendAsync(ctx, func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			fmt.Printf("!!!! Receive message error: %s\n", err)
		} else {
			// redis.Redis.HSET(uuid.NewV4().String(), uuid.NewV4().String())
			fmt.Printf("!!! Send message success: result= %s\n", result.String())
		}
	}, sendData...)
	if err != nil {
		fmt.Printf("!!! Send data error: %s", err.Error())
		os.Exit(1) // 进行日志记录或者警告，这里直接退出程序
	}

	// CloseProducer() // 生产者记得及时关闭, 因为这里定义了全局变量，程序结束时回收资源
	return data, nil
}

/*
// 修改值必须是指针类型否则不可行
	//var sendData []*primitive.Message // 将多条消息整合发送
	//if v.Kind() == reflect.Slice {
	//	l := v.Len()
	//	for i := 0; i < l; i++ {
	//		value := v.Index(i) // Value of item
	//		infoByte, _ := json.Marshal(value.Interface())
	//
	//		msg := &primitive.Message{
	//			Topic: _const.Topic,
	//			Body:  infoByte,
	//			Queue: &primitive.MessageQueue{
	//				Topic:      _const.Topic,
	//				BrokerName: _const.Broker,
	//				QueueId:    0,
	//			},
	//		}
	//		sendData = append(sendData, msg)
	//	}
	//} else {
	//	//记录日志 ...
	//	infoByte, _ := json.Marshal(data)
	//	pmsg := &primitive.Message{
	//		Topic: _const.Topic,
	//		Body:  infoByte,
	//		Queue: &primitive.MessageQueue{ // 向指定 queue 发送信息
	//			QueueId:    0,
	//			Topic:      _const.Topic,
	//			BrokerName: _const.BrokerName,
	//		},
	//	}
	//	sendData = append(sendData, pmsg)
	//}

	//msg := []*primitive.Message{&primitive.Message{
	//	Topic: _const.Topic,
	//	Body:  []byte("this is a message body1"),
	//	Queue: &primitive.MessageQueue{
	//		Topic:      _const.Topic,
	//		BrokerName: _const.BrokerName,
	//		QueueId:    0,
	//	},
	//}, &primitive.Message{
	//	Topic: _const.Topic,
	//	Body:  []byte("this is a message body2"),
	//	Queue: &primitive.MessageQueue{
	//		Topic:      _const.Topic,
	//		BrokerName: _const.BrokerName,
	//		QueueId:    0,
	//	},
	//}}
*/
