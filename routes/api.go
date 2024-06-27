package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-tool-api/app/controller"
	"net/http"
)

var Router *gin.Engine

var Api = new(apiRouter)

type apiRouter struct{}

func (apiRouter *apiRouter) InitRouter() {
	Router = gin.Default()
	//_ = Router.SetTrustedProxies([]string{"127.0.0.1"})
	//跨域中间件
	Router.Use(cors())
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

		printer := api.Group("printer")
		{
			printer.GET("list", controller.Print.List)
			printer.GET("printFile", controller.Print.PrintFile)
		}

	}

	if err := Router.Run(":8080"); err != nil {
		fmt.Printf("error:" + err.Error())
	}
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
