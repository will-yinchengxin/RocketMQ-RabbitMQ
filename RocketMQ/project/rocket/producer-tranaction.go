package main


//事务监听器
type Trans struct {
	localTrans       *sync.Map
	transactionIndex int32
}

func NewTrans() *Trans {
	return &Trans{
		localTrans: new(sync.Map),
	}
}

func (tr *Trans) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	switch msg.GetTags() {
	case "taga":
		return primitive.CommitMessageState
	case "tagb":
		return primitive.RollbackMessageState
	case "tagc","tagd":
		return primitive.UnknowState
	}
	return primitive.UnknowState
}

func (tr *Trans) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	switch msg.GetTags() {
	case "tagc":
		log.Printf("check commit : %v\n", string(msg.Body) )
		return primitive.CommitMessageState
	case "tagd":
		log.Printf("check rollback: %v\n", string(msg.Body))
		return primitive.RollbackMessageState
	}
	return primitive.RollbackMessageState
}

func main() {
	p, _ := rocketmq.NewTransactionProducer(
		NewTrans(),
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(1),
	)
	err := p.Start()
	defer p.Shutdown()
	if err != nil {
		log.Printf("start producer error: %s\n", err.Error())
		os.Exit(1)
	}
	tags := [4]string{"taga","tagb","tagc","tagd"}
	topic := "trans"
	for i := 0; i < 4; i++ {
		msg := primitive.NewMessage(topic, []byte("Hello RocketMQ"+strconv.Itoa(i)))
		msg.WithTag(tags[i])
		res, err := p.SendMessageInTransaction(context.Background(), msg)
		if err != nil {
			log.Printf(err.Error())
			return
		} else {
			log.Printf("预投递成功 resid : %s \n",res.MsgID)
		}
	}
	time.Sleep(20 * time.Minute)
}
/**
*
2022/01/25 12:57:00 预投递成功 resid : C0A80165CBC0000000007e5ffd600001
2022/01/25 12:57:00 预投递成功 resid : C0A80165CBC0000000007e5ffd600002
2022/01/25 12:57:00 预投递成功 resid : C0A80165CBC0000000007e5ffd600003
2022/01/25 12:57:00 预投递成功 resid : C0A80165CBC0000000007e5ffd600004
2022/01/25 12:57:07 check rollback: Hello RocketMQ3
2022/01/25 12:57:07 check commit : Hello RocketMQ2
*/
