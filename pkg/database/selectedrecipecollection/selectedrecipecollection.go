package selectedrecipecollection

import (
	"log"
	"masterchef_bot/pkg/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SelectedRecipe document
type SelectedRecipe struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserID   primitive.ObjectID `bson:"UserID"`
	ChatID   int64              `bson:"ChatID"`
	Name     string             `bson:"Name"`
	URL      string             `bson:"URL"`
	Selected time.Time          `bson:"Selected"`
}

const collection = "selectedrecipes"

// Save new selection for user in a given chat
func Save(name string, url string, chatID int64, userID primitive.ObjectID) (err error) {

	userFilter := bson.M{
		"ChatID": chatID,
		"UserID": userID,
	}
	update := bson.M{
		"$set": bson.M{
			"Name":     name,
			"URL":      url,
			"Selected": time.Now(),
		},
	}

	upsert := true
	opt := options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	}
	// Try to update if the exsiting user and chat is found
	res := database.Manager.Get(collection).FindOneAndUpdate(*database.Manager.GetContext(), userFilter, update, &opt)
	if err = res.Err(); err != nil {
		log.Print(res.Err())
	}
	return
}

// GetByUser from the collection
func GetByUser(userID primitive.ObjectID) (recipe *SelectedRecipe) {

	recipe = &SelectedRecipe{}
	filter := bson.M{
		"UserID": userID,
	}

	opt := options.FindOne()
	// Sort by `Selected` field descending
	opt.SetSort(bson.D{
		bson.E{
			Key: "Selected", Value: -1,
		},
	})

	err := database.Manager.Get(collection).FindOne(*database.Manager.GetContext(), filter, opt).Decode(recipe)
	if err != nil {
		log.Print(err)
		recipe = nil
	}

	return
}
