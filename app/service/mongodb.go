package service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
	"time"
)

var Mongodb = new(MongodbServer)

type MongodbServer struct {
}

type SpuParams struct {
	Company   string `json:"company"`
	UniCommID string `json:"uniCommID"`
}
type SpuModel struct {
	ID        string `json:"id"`
	AddAt     string `json:"addAt"`
	AddBy     string `json:"addBy"`
	CommCode  string `json:"commCode"`
	UniCommID string `json:"uniCommID"`
}
type Params struct {
	Url        string `json:"url"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
	Filter     any    `json:"filter"`
	Options    any    `json:"options"`
}

func (server MongodbServer) Find(param Params) ([]any, error) {
	/// 检测 数据
	param, err := server.CheckParamData(param)
	if err != nil {
		return nil, err
	}

	// 设置 MongoDB 连接选项
	clientOptions := options.Client().ApplyURI(param.Url).SetConnectTimeout(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 连接 MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Failed to connect to MongoDB:" + err.Error())
		return nil, errors.New("Failed to connect to MongoDB:" + err.Error())
	}
	// 确保在最后关闭链接
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			fmt.Println("关闭连接失败！" + err.Error())
		}
		fmt.Println("已经和Mongodb数据库断开链接！")
	}()

	// 定义查询条件
	// 获取数据库引用
	collection := client.Database(param.Database).Collection(param.Collection)
	// 执行查询
	cursor, err := collection.Find(ctx, param.Filter)
	if err != nil {
		log.Println("查询失败1:" + err.Error())
		return nil, errors.New("查询失败！" + err.Error())
	}
	// 遍历游标并获取结果
	var results []any
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Println("查询失败2:" + err.Error())
		return nil, errors.New("make data err ！" + err.Error())
	}
	return results, nil
}

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (server MongodbServer) FindOne(param Params) (any, error) {
	/// 检测 数据
	param, err := server.CheckParamData(param)
	if err != nil {
		return nil, err
	}
	// 设置 MongoDB 连接选项
	clientOptions := options.Client().ApplyURI(param.Url).SetConnectTimeout(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 连接 MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.New("Failed to connect to MongoDB:" + err.Error())
	}
	// 确保在最后关闭链接
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			fmt.Println("关闭连接失败！" + err.Error())
		}
	}()
	// 获取数据库引用
	collection := client.Database(param.Database).Collection(param.Collection)
	// 执行查询
	//// 遍历游标并获取结果
	//var results any
	results := struct {
		Id string `json:"id"`
	}{}
	//var results = [][]

	err = collection.FindOne(ctx, param.Filter).Decode(&results)
	//err = collection.FindOne(ctx, param.Filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, errors.New("Raw err:" + err.Error())
	}

	fmt.Printf(" %v - %T ", results, results)

	return results, nil
}

func (server MongodbServer) CheckParamData(param Params) (Params, error) {
	if param.Url == "" {
		return param, errors.New("url 不能为空！")
	}
	if param.Database == "" {
		return param, errors.New("database 不能为空！")
	}
	if param.Collection == "" {
		return param, errors.New("collection 不能为空！")
	}
	return param, nil
}

func isArray(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Array
}
