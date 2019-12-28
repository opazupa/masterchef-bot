package usercollection

import (
	"masterchef_bot/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User document
type User struct {
	DatabaseID primitive.ObjectID `bson:"_id"`
	ID         int
	Name       string
}

const user = "user"

// Create new user
func Create(id int, name string) *User {
	newUser := bson.M{
		"ID":   id,
		"Name": name,
	}

	var result User

	db := database.Get()
	_, err := db.Collection(user).InsertOne(database.GetContext(), newUser)
	if err != nil {
		return nil
	}
	return &result
}

// Get user by id
func Get(id int) *User {
	filter := bson.D{
		primitive.E{
			Key: "ID", Value: id,
		},
	}

	var result User

	db := database.Get()
	err := db.Collection(user).FindOne(database.GetContext(), filter).Decode(&result)
	if err != nil {
		return nil
	}

	return &result
}
