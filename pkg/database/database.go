package database

import (
	"context"
	"log"
	"masterchef_bot/pkg/configuration"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Get mongo database
func Get() *mongo.Database {
	configuration := configuration.Get()
	clientOptions := options.Client().ApplyURI(configuration.DatabaseConnection)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Panic(err)
	}

	// Check connections
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Printf("Connected to the [%s] database.", configuration.DatabaseName)
	}

	return client.Database(configuration.DatabaseName)
}
