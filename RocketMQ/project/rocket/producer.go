package rocket

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"reflect"
	_const "rock/const"
	"rock/dto"
)

type Producer struct {
}

func (p *Producer) Send() (err error){
	data := dto.GetData()
	v := reflect.ValueOf(data)

	// 修改值必须是指针类型否则不可行
	var sendData []*primitive.Message
	if v.Kind() == reflect.Slice {
		l := v.Len()
		for i := 0; i < l; i++ {
			value := v.Index(i) // Value of item
			infoByte, _ := json.Marshal(value.Interface())
			pmsg := &primitive.Message{
				Topic: _const.Topic,
				Body:  infoByte,
			}
			sendData = append(sendData, pmsg)
		}

	} else {
		//记录日志 ...
		infoByte, _ := json.Marshal(data)
		pmsg := &primitive.Message{
			Topic: _const.Topic,
			Body:  infoByte,
		}
		sendData = append(sendData, pmsg)
	}
	ctx := context.TODO() // 这里可以做成外部传递的形式
	err = RocketmqProducerClient.SendAsync(ctx, func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			return
		}
	}, sendData...)
	return nil
}
