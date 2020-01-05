package usercollection

import (
	"masterchef_bot/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User document
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserName string             `bson:"UserName"`
}

const userCollection = "users"

// Create new user
func Create(userName string) (*User, error) {
	newUser := bson.M{
		"UserName": userName,
	}

	db := database.Get()
	inserted, err := db.Collection(userCollection).InsertOne(database.GetContext(), newUser)
	if err != nil {
		return nil, err
	}
	return getByID(inserted.InsertedID.(primitive.ObjectID)), nil
}

// GetByID from collection
func getByID(id primitive.ObjectID) *User {

	filter := bson.D{
		primitive.E{
			Key: "_id", Value: id,
		},
	}
	var result User

	db := database.Get()
	err := db.Collection(userCollection).FindOne(database.GetContext(), filter).Decode(&result)
	if err != nil {
		return nil
	}

	return &result
}

// GetByUserName from collection
func GetByUserName(userName *string) *User {

	if userName == nil {
		return nil
	}

	filter := bson.D{
		primitive.E{
			Key: "UserName", Value: userName,
		},
	}
	var result User

	db := database.Get()
	err := db.Collection(userCollection).FindOne(database.GetContext(), filter).Decode(&result)
	if err != nil {
		return nil
	}

	return &result
}
