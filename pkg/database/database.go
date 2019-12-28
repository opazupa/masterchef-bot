package database

import (
	"context"
	"log"
	"masterchef_bot/pkg/configuration"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Check mongo database
func Check() {
	configuration := configuration.Get()
	clientOptions := options.Client().ApplyURI(configuration.DatabaseConnection)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Panic(err)
	}

	ctx := GetContext()

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
}

// Get mongo database
func Get() *mongo.Database {
	configuration := configuration.Get()
	clientOptions := options.Client().ApplyURI(configuration.DatabaseConnection)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Panic(err)
	}

	ctx := GetContext()
	// Check connections
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(configuration.DatabaseName)
}

// GetContext for database call
func GetContext() context.Context {
	return context.Background()
}
