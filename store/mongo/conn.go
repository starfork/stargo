package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Config struct {
	Password, User, DbName, Addr string
	Monitor                      bool
	//FileServer fileServer
}

func NewConnection(c Config) {
	auth := options.Credential{
		//AuthSource: "<authenticationDb>",
		Username: c.User,
		Password: c.Password,
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + c.Addr).SetAuth(auth)
	//debug
	if c.Monitor {
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				log.Println(evt.Command.String())
				log.Println(evt.CommandName)
			},
		}
		clientOptions.SetMonitor(cmdMonitor)
	}
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Connected to MongoDB!")
}

func Db() *mongo.Client {
	return client
}
