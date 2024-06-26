package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redis-tool/app/models/response"
	"redis-tool/app/service"
)

var Mongodb = new(MongodbController)

type MongodbController struct{}

func (t *MongodbController) Find(c *gin.Context) {
	var param service.Params
	// 解析 JSON 数据
	if err := c.ShouldBindJSON(&param); err != nil {
		// 处理解析错误
		c.JSON(http.StatusOK, response.Fail("数据格式异常 ！"+err.Error()))
		return
	}
	data, err := service.Mongodb.Find(param)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Ok(data))
	return
}

func (t *MongodbController) FindOne(c *gin.Context) {
	var param service.Params
	// 解析 JSON 数据
	if err := c.ShouldBindJSON(&param); err != nil {
		// 处理解析错误
		c.JSON(http.StatusOK, response.Fail("数据格式异常 ！"+err.Error()))
		return
	}

	data, err := service.Mongodb.FindOne(param)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Ok(data))
	return
}
