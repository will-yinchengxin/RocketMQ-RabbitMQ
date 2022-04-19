package rocket

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
	"os"
	_const "rock/const"
)

var RocketmqProducerClient rocketmq.Producer
var RocketmqPushConsumerClient rocketmq.PushConsumer
var RocketmqPullConsumerClient rocketmq.PullConsumer

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
		GroupName: _const.BrokerName,
		Topic:     _const.Topic,
	}
}

func InitRocket() {
	var err error
	rocket := NewRocketmqConfig()

	// 生产者部分
	RocketmqProducerClient, err = rocketmq.NewProducer(
		producer.WithNameServer(rocket.Host),
		producer.WithRetry(rocket.Retry),
		producer.WithGroupName(rocket.GroupName),
		//producer.WithSendMsgTimeout(10*time.Second),  // 默认是 3 秒中
		/*
			设置 queue 的选择规则：
				NewManualQueueSelector： 指定 queue 发送
					msg := &primitive.Message{
						Topic: _const.Topic,
						Body:  infoByte,
						Queue: &primitive.MessageQueue{  // 通过这里进行设置特定的 QueueId
							Topic:      _const.Topic,
							BrokerName: _const.BrokerName,
							QueueId:    0,
						},
					}


				NewRandomQueueSelector：先创建randomQueueSelector，然后设置其rander；Select方法通过 r.rander.Intn(len(queues))随机选择index，然后从queue取值
				(默认方式)NewRoundRobinQueueSelector：qIndex为： int(i) % len(queues) 的结果， 然后负载均衡
				NewHashQueueSelector：qIndex为： int(hasher.Sum32()) % len(queues) 的结果
		*/
		//producer.WithQueueSelector(producer.NewManualQueueSelector()),
	)
	if err != nil {
		fmt.Printf("new producer error: %s", err.Error())
		os.Exit(1)
	}
	log.Print("new producer success")
	err = RocketmqProducerClient.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	log.Print("start producer success")

	/*
		消费者部分
				消费者有两中获取消息的方式，本质都是 pull
				- pull
				- push
	*/
	// push，底层是 pull 模式进行封装的，不需要客户端管理消费进度
	RocketmqPushConsumerClient, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer(rocket.Host),
		consumer.WithConsumerModel(consumer.Clustering), // 设置接受消息的模式 -BroadCasting： 广播模式    -Clustering：集群模式
		consumer.WithGroupName("rocketTest"),

		//consumer.WithConsumerOrder(true), // 设置消费者局部消费有序
	)
	if err != nil {
		fmt.Printf("new push consumer error: %s", err.Error())
		os.Exit(1)
	}
	log.Print("new push consumer success")

	// pull 需要自行管理消费的进度 (目前官方的 api 还未支持)
	//RocketmqPullConsumerClient, err = rocketmq.NewPullConsumer(
	//	consumer.WithGroupName("rocketTestTwo"),
	//	consumer.WithNsResolver(primitive.NewPassthroughResolver(rocket.Host)),
	//)
	//if err != nil {
	//	fmt.Printf("new pull consumer error: %s", err.Error())
	//	os.Exit(1)
	//}
	//log.Print("new pull Consumer success")
	//err = RocketmqPullConsumerClient.Start()
	//if err != nil {
	//	fmt.Printf("start pull consumer error: %s", err.Error())
	//	os.Exit(1)
	//}
	//log.Print("start pull consumer success")
}

// 创建 topic
func CreateTopic() {
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver([]string{_const.NameServer})))
	err = testAdmin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate(_const.Topic),
		admin.WithWriteQueueNums(4), // 指定创建 queue 的数量，queue 可细分为 WriteQueue 与 ReadQueue 有映射关系
		admin.WithReadQueueNums(4),  // 所以创建 WriteQueue 与 ReadQueue 的数量要相等
		admin.WithBrokerAddrCreate(_const.Broker),
		admin.WithPerm(6), // 设置该 Topic 的读写模式 6：同时支持读写 4：禁写 2：禁读
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

func CloseProducer() {
	err := RocketmqProducerClient.Shutdown()
	if err != nil {
		fmt.Printf("CloseProducer error: %s", err.Error())
	}
	log.Println("CloseProducer Success")
}
func ClosePushConsumer() {
	err := RocketmqPushConsumerClient.Shutdown()
	if err != nil {
		fmt.Printf("ClosePushConsumer error: %s", err.Error())
	}
	log.Println("ClosePushConsumer Success")
}
func ClosePullConsumer() {
	RocketmqPullConsumerClient.Shutdown()

}
