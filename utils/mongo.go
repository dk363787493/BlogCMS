package utils

import (
	"BlogCMS/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var mongoClient *mongo.Client

func init() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(config.Configuration.Mongo.Url)

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	mongoClient = client
	if err != nil {
		log.Fatal(err)
	}
}

func GetMongoClient() *mongo.Client {
	return mongoClient
}
