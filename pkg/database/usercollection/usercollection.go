package usercollection

import (
	"log"
	"masterchef_bot/pkg/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User document
type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserName   string             `bson:"UserName"`
	TelegramID int                `bson:"TelegramID"`
	Registered time.Time          `bson:"Registered"`
}

const collection = "users"

// Create new user
func Create(userName string, id int) (*User, error) {
	newUser := bson.M{
		"UserName":   userName,
		"TelegramId": id,
		"Registered": time.Now(),
	}

	inserted, err := database.Manager.Get(collection).InsertOne(*database.Manager.GetContext(), newUser)
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

	err := database.Manager.Get(collection).FindOne(*database.Manager.GetContext(), filter).Decode(&result)
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

	err := database.Manager.Get(collection).FindOne(*database.Manager.GetContext(), filter).Decode(&result)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &result
}
