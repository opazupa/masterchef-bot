package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User document
type User struct {
	DatabaseID primitive.ObjectID `bson:"_id"`
	ID         int
	Name       string
}

const user = "user"

// Create new user
func Create(db *mongo.Database, id int, name string) *User {
	newUser := bson.M{
		"ID":   id,
		"Name": name,
	}

	var result User

	_, err := db.Collection(user).InsertOne(context.Background(), newUser)
	if err != nil {
		return nil
	}
	return &result
}

// Get user by id
func Get(db *mongo.Database, id int) *User {
	filter := bson.D{
		primitive.E{
			Key: "ID", Value: id,
		},
	}

	var result User

	db.Collection(user).FindOne(context.Background(), filter).Decode(&result)
	return &result
}
