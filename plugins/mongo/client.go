package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Mongo interface {
	Init()
}

type mongoClient struct {
	client *mongo.Client
}

func (m *mongoClient) Init() {

	ctx,_:=context.WithTimeout(context.TODO(),time.Second*500)
	err:=m.client.Connect(ctx)
	if err!=nil{
		panic(err)
	}
	err=m.client.Ping(ctx,readpref.Primary())
	if err!=nil{
		panic(err)
	}
}

func NewMongo() Mongo{
	option:=options.Client()
	option=option.SetMinPoolSize(1)
	option=option.ApplyURI("mongodb://root:pwd@localhost:27017")
	client,err:=mongo.NewClient(option)
	if err!=nil{
		panic(err)
	}
	return &mongoClient{client:client}
}