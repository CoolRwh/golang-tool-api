package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"redis-tool/app/controller"
)

func main() {
	InitRouter()
}

var Router *gin.Engine

func InitRouter() {
	Router = gin.Default()
	api := Router.Group("api/v1")
	{
		redis := api.Group("redis")
		{
			redis.GET("config", controller.Redis.Config)
			redis.GET("keys", controller.Redis.Keys)
			redis.GET("info", controller.Redis.Info)
		}

		mongodb := api.Group("mongodb")
		{
			mongodb.POST("find", controller.Mongodb.Find)
			mongodb.POST("fineOne", controller.Mongodb.FindOne)
		}
	}
	if err := Router.Run(":8080"); err != nil {
		fmt.Printf("error:" + err.Error())
	}
}
