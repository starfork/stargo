package mongo

import (
	"context"
	"log"

	"github.com/starfork/stargo/store"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client *mongo.Client
	c      *store.Config
	ctx    context.Context
}

func NewMongo(config *store.Config) store.Store {
	return &Mongo{
		c:   config,
		ctx: context.Background(),
	}
}

func (e *Mongo) Connect(conf ...*store.Config) {
	c := e.c
	auth := options.Credential{
		//AuthSource: "<authenticationDb>",
		Username: c.User,
		Password: c.Auth,
	}
	var client *mongo.Client
	clientOptions := options.Client().ApplyURI("mongodb://" + c.Host).SetAuth(auth)
	//debug
	if c.Monitor {
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				//log.Println(evt.Command.String())
				//log.Println(evt.CommandName)
			},
		}
		clientOptions.SetMonitor(cmdMonitor)
	}
	var err error
	client, err = mongo.Connect(e.ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = client.Ping(e.ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	e.client = client

}

func (e *Mongo) GetInstance(conf ...*store.Config) *mongo.Client {
	if len(conf) > 0 {
		e.Connect(conf...)
		return e.client
	}
	return e.client
}

func (e *Mongo) Close() {
	e.client.Disconnect(e.ctx)
}
