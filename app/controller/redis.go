package controller

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang-tool-api/app/models/response"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var Redis = new(RedisController)

type RedisController struct{}

type Config struct {
	Addr     string `json:"addr"`
	Password string `json:"pwd"`
	DB       int    `json:"DB"`
	Key      string `json:"key"`
	Start    int64  `json:"start"`
	Stop     int64  `json:"stop"`
}

var ctx = context.Background()

func (t RedisController) Config(c *gin.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.56.56:6379", // Redis地址
		Password: "123456",             // Redis密码，如果没有则为空字符串
		DB:       0,                    // 使用默认DB
	})
	config, err := rdb.Info(ctx, "all").Result()
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	configNew, err := parseInfo(config)
	if err != nil {
		errMsg := &response.BusinessError{Code: response.RequestParamError, Msg: err.Error()}
		c.JSON(http.StatusOK, response.ResultCustom(errMsg))
		return
	}
	// 关闭连接
	err = rdb.Close()
	if err != nil {
		errMsg := &response.BusinessError{Code: response.RedisClientCloseError, Msg: err.Error()}
		c.JSON(http.StatusOK, response.ResultCustom(errMsg))
		return
	}
	c.JSON(http.StatusOK, response.Ok(configNew))
	return
}

func (t RedisController) Keys(c *gin.Context) {
	var configData Config
	err := c.ShouldBindJSON(&configData)
	if err != nil {
		errMsg := &response.BusinessError{Code: response.RequestParamError, Msg: err.Error()}
		c.JSON(http.StatusOK, response.ResultCustom(errMsg))
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     configData.Addr,     // Redis地址
		Password: configData.Password, // Redis密码，如果没有则为空字符串
		DB:       configData.DB,       // 使用默认DB
	})
	if configData.Key == "" {
		configData.Key = "*"
	}
	list, _ := rdb.Keys(ctx, configData.Key).Result()
	// 关闭连接
	err = rdb.Close()
	if err != nil {
		errMsg := &response.BusinessError{Code: response.RequestParamError, Msg: err.Error()}
		c.JSON(http.StatusOK, response.ResultCustom(errMsg))
		return
	}
	c.JSON(200, gin.H{"list": list, "total": len(list)})
	return
}

func (t RedisController) Info(c *gin.Context) {
	var configData Config
	err := c.ShouldBindJSON(&configData)
	if err != nil {
		c.JSON(http.StatusOK, response.ResultCustom(&response.BusinessError{Code: response.RequestParamError, Msg: err.Error()}))
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     configData.Addr,     // Redis地址
		Password: configData.Password, // Redis密码，如果没有则为空字符串
		DB:       configData.DB,       // 使用默认DB
	})

	type result struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		TTL   int    `json:"TTL"`
		Size  int    `json:"size"`
		Data  any    `json:"data"`
		Info  any    `json:"info"`
		Total int    `json:"total"`
	}
	var resultData result
	resultData.Name = configData.Key
	resultData.Type, _ = rdb.Type(ctx, configData.Key).Result()
	resultData.Info, _ = rdb.DebugObject(ctx, configData.Key).Result()
	switch resultData.Type {
	case "string":
		resultData.Data, _ = rdb.Get(ctx, configData.Key).Result()
		break
	case "list":
		resultData.Data, _ = rdb.LRange(ctx, configData.Key, configData.Start, configData.Stop).Result()
		break
	case "set":
		resultData.Data, _ = rdb.SMembers(ctx, configData.Key).Result()
		break
	case "zset":
		resultData.Data, _ = rdb.ZRangeWithScores(ctx, configData.Key, configData.Start, configData.Stop).Result()
		break
	case "hash":
		resultData.Data, _ = rdb.HGetAll(ctx, configData.Key).Result()
		break
	case "stream":
		break
	}
	// 关闭连接
	err = rdb.Close()
	if err != nil {
		c.JSON(http.StatusOK, response.ResultCustom(&response.BusinessError{Code: response.RedisClientError, Msg: err.Error()}))
		return
	}
	c.JSON(http.StatusOK, response.Ok(resultData))
	return
}

// parseInfoLine 解析 INFO 命令返回的单行信息
func parseInfoLine(line string) (key string, value interface{}, err error) {
	// 假设行格式为 "key:value" 或 "key=value"
	parts := regexp.MustCompile(`\s*:\s*|\s*=\s*`).Split(line, 2)
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("invalid INFO line format: %s", line)
	}
	key = strings.ToLower(parts[0]) // 将键转换为小写以进行标准化
	valueStr := parts[1]

	// 尝试解析为数字（整数或浮点数），如果失败则保留为字符串
	if _, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		// 整数
		intValue, _ := strconv.ParseInt(valueStr, 10, 64)
		return key, intValue, nil
	}
	if _, err := strconv.ParseFloat(valueStr, 64); err == nil {
		// 浮点数
		floatValue, _ := strconv.ParseFloat(valueStr, 64)
		return key, floatValue, nil
	}

	// 默认作为字符串处理
	return key, valueStr, nil
}

// parseInfo 解析 INFO 命令的输出到 map[string]interface{}
func parseInfo(info string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	scanner := bufio.NewScanner(strings.NewReader(info))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line) // 去除行首尾的空白字符
		if strings.HasPrefix(line, "#") {
			// 忽略注释行
			continue
		}
		// 忽略 空格行
		if len(line) <= 0 {
			continue
		}
		key, value, err := parseInfoLine(line)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
