package main

import (
	"fmt"
	"rock/rocket"
	"time"
)

func main() {
	// 初始化连接
	rocket.InitRocket()
	defer rocket.Close()
	// 生产者
	//pro := rocket.Producer{}
	//pro.Send()


	// 消费者
	for {
		con := rocket.Consumer{}
		con.Start()
		time.Sleep(time.Second*2)
		fmt.Println("get message!")
	}
}
