package database

import (
	"context"
	"log"
	"masterchef_bot/pkg/configuration"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoManager interface
type MongoManager interface {
	Get(collection string) *mongo.Collection
	GetContext() *context.Context
}

type manager struct {
	db  *mongo.Database
	ctx *context.Context
}

// Manager for mongo db
var Manager MongoManager

func init() {

	configuration := configuration.Get()
	clientOptions := options.Client().ApplyURI(configuration.DatabaseConnection)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()

	// Check connections
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Printf("Connected to the [%s] database.", configuration.DatabaseName)
	}
	Manager = &manager{db: client.Database(configuration.DatabaseName), ctx: &ctx}
}

func (mgr *manager) Get(collectionName string) *mongo.Collection {
	return mgr.db.Collection(collectionName)
}

func (mgr *manager) GetContext() *context.Context {
	return mgr.ctx
}
