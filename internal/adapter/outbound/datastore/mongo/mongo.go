package mongo

import (
	"context"
	"sync"
	"time"

	"github.com/ilmimris/poc-go-ulid/pkg/database"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
)

var once sync.Once
var mongodb *mongo.Database
var client *database.Database

type OptMongo struct {
	URI               string
	DB                string
	AppName           string
	ConnectionTimeOut int
	PingTimeOut       int
}

func New(configs ...OptMongo) *mongo.Database {
	var cfg OptMongo
	if len(configs) > 0 {
		cfg = configs[0]
	}

	if mongodb != nil {
		return mongodb
	}

	once.Do(func() {
		c := &database.Client{
			URI:            cfg.URI,
			DB:             cfg.DB,
			AppName:        cfg.AppName,
			ConnectTimeout: time.Duration(cfg.ConnectionTimeOut) * time.Second,
			PingTimeout:    time.Duration(cfg.PingTimeOut) * time.Second,
		}
		client = database.MongoConnectClient(c)
		mongodb = client.Database
		log.Infof("MongoDB connected")
	})

	return mongodb

}

func GetMongoDB() *mongo.Database {
	return mongodb
}

func Close(ctx context.Context) {
	if client != nil {
		err := client.Database.Client().Disconnect(ctx)
		if err != nil {
			log.Error(err)
		}
		log.Info("MongoDB closed")
	}
}
