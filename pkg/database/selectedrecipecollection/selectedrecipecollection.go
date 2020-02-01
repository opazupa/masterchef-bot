package selectedrecipecollection

import (
	"masterchef_bot/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	userFilter := bson.M{
		"ChatID": chatID,
		"UserID": userID,
	}
	update := bson.M{
		"$set": bson.M{
			"Name": name,
			"URL":  url,
		},
	}

	upsert := true
	opt := options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	}
	// Try to update if the exsiting user and chat is found
	res := database.Manager.Get(collection).FindOneAndUpdate(*database.Manager.GetContext(), userFilter, update, &opt)

	return res.Err()
}

// GetByUser from the collection
func GetByUser(userID primitive.ObjectID) *SelectedRecipe {

	filter := bson.M{
		"UserID": userID,
	}

	var result SelectedRecipe

	err := database.Manager.Get(collection).FindOne(*database.Manager.GetContext(), filter).Decode(&result)
	if err != nil {
		return nil
	}

	return &result
}
