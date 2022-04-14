package rocket

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
	_const "rock/const"
)

var RocketmqProducerClient rocketmq.Producer
var RocketmqConsumerClient rocketmq.PushConsumer

// 配置文件可以配置在 yml 文件中, 通过 viper 进行读取, 这里方便起见直接进行配置
type RocketmqConfig struct {
	Host      []string `yaml:"host,omitempty"`
	Retry     int      `yaml:"retry,omitempty"`
	GroupName string   `yaml:"groupName,omitempty"`
	Topic     string   `yaml:"topic,omitempty"`
}

func NewRocketmqConfig() *RocketmqConfig {
	return &RocketmqConfig{
		Host:      []string{_const.NameServer},
		Retry:     2,
		GroupName: _const.GroupName,
		Topic:     _const.Topic,
	}
}

// 创建 topic
func CreateTopic() {
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver([]string{_const.NameServer})))
	err = testAdmin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate("newTopic"),
		admin.WithBrokerAddrCreate(_const.Broker),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// 删除 topic
func DeleteTopic() {
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver([]string{_const.NameServer})))
	err = testAdmin.DeleteTopic(
		context.Background(),
		admin.WithTopicDelete("newTopic"),
		admin.WithBrokerAddrDelete(_const.Broker),
		//admin.WithNameSrvAddr(nameSrvAddr),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func InitRocket() {
	var err error
	rocket := NewRocketmqConfig()
	RocketmqProducerClient, err = rocketmq.NewProducer(
		producer.WithNameServer(rocket.Host),
		producer.WithRetry(rocket.Retry),
		producer.WithGroupName(rocket.GroupName),
	)
	if err != nil {
		panic(err)
	}

	// 生产者
	err = RocketmqProducerClient.Start()
	if err != nil {
		panic(err)
	}

	// 消费者
	RocketmqConsumerClient, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer(rocket.Host),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(rocket.GroupName),
	)
	if err != nil {
		panic(err)
	}
}

func Close() {
	RocketmqProducerClient.Shutdown()
	RocketmqConsumerClient.Shutdown()
}
