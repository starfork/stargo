package mongo

import (
	"context"
	"log"

	"github.com/starfork/stargo/config"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client *mongo.Client
}

func Connect(conf *config.Config) *Mongo {
	c := conf.MongoDb
	auth := options.Credential{
		//AuthSource: "<authenticationDb>",
		Username: c.User,
		Password: c.Password,
	}
	var client *mongo.Client
	clientOptions := options.Client().ApplyURI("mongodb://" + c.Host).SetAuth(auth)
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
	return &Mongo{
		client: client,
	}
	//fmt.Println("Connected to MongoDB!")
}

func (e *Mongo) GetInstance(conf ...*config.Config) *mongo.Client {
	if len(conf) > 0 {
		rs := Connect(conf[0])
		return rs.client
	}
	return e.client
}

func (e *Mongo) Close() {

}
