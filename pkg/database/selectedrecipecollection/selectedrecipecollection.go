package selectedrecipecollection

import (
	"masterchef_bot/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SelectedRecipe document
type SelectedRecipe struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID primitive.ObjectID `bson:"UserID"`
	ChatID int64              `bson:"ChatID"`
	Name   string             `bson:"Name"`
	URL    string             `bson:"URL"`
}

const collection = "selectedrecipes"

// Save new selection for user in a given chat
func Save(name string, url string, chatID int64, userID primitive.ObjectID) error {
	newUserSelection := bson.M{
		"Name":   name,
		"URL":    url,
		"ChatID": chatID,
		"UserID": userID,
	}

	_, err := database.Manager.Get(collection).InsertOne(*database.Manager.GetContext(), newUserSelection)
	return err
}

// GetByUser for given chat
func GetByUser(chatID int64, userID primitive.ObjectID) *SelectedRecipe {

	filter := bson.M{
		"ChatID": chatID,
		"UserID": userID,
	}

	var result SelectedRecipe

	err := database.Manager.Get(collection).FindOne(*database.Manager.GetContext(), filter).Decode(&result)
	if err != nil {
		return nil
	}

	return &result
}
