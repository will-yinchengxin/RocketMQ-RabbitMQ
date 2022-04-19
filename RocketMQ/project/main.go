package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rock/core"
	"rock/redis"
	"rock/rocket"
)

func init() {
	rocket.InitRocket()
	redis.InitRedis()
}

func main() {
	r := gin.Default()

	r.GET("/rocket", func(c *gin.Context) {
		// 初始化连接
		rocket.CreateTopic()
		c.String(http.StatusOK, "init rocket")
	})
	r.GET("/send", func(c *gin.Context) {
		// 初始化连接
		pro := rocket.Producer{}
		res, _ := pro.Send()

		c.JSON(http.StatusOK, res)
	})
	r.GET("/consume", func(c *gin.Context) {
		// 初始化连接
		pro := rocket.Consumer{}
		pro.PushConsumerStart()
		//pro.PullConsumerStart()
		c.JSON(http.StatusOK, "get msg success")
	})

	defer func() {
		core.Stop()
	}()

	err := endless.ListenAndServe(":8050", r)
	if err != nil {
		core.Stop()
		log.Fatalf("sever: %s\n", err)
	}
	log.Println("logout success")
}
