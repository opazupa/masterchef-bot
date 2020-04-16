package usercollection

import (
	"log"
	"masterchef_bot/pkg/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Favourite
type favourite struct {
	RecipeID primitive.ObjectID `bson:"RecipeID"`
	Added    time.Time          `bson:"Added"`
}

// User document
type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserName   string             `bson:"UserName"`
	TelegramID int                `bson:"TelegramID"`
	Registered time.Time          `bson:"Registered"`
	Favourites []favourite        `bson:"Favourites"`
}

const collection = "users"

// Create new user
func Create(userName string, id int) (user *User, err error) {
	newUser := bson.M{
		"UserName":   userName,
		"TelegramId": id,
		"Registered": time.Now(),
		"Favourites": []favourite{},
	}

	inserted, err := database.Manager.Get(collection).InsertOne(*database.Manager.GetContext(), newUser)
	if err != nil {
		return
	}
	user, err = getByID(inserted.InsertedID.(primitive.ObjectID))
	return
}

// GetByID from collection
func getByID(id primitive.ObjectID) (user *User, err error) {

	user = &User{}
	filter := bson.D{
		primitive.E{
			Key: "_id", Value: id,
		},
	}
	err = database.Manager.Get(collection).FindOne(*database.Manager.GetContext(), filter).Decode(user)
	if err != nil {
		log.Print(err)
		user = nil
	}
	return
}

// GetByUserName from collection
func GetByUserName(userName *string) (user *User) {

	if userName == nil {
		return
	}

	user = &User{}
	filter := bson.D{
		primitive.E{
			Key: "UserName", Value: userName,
		},
	}

	err := database.Manager.Get(collection).FindOne(*database.Manager.GetContext(), filter).Decode(user)
	if err != nil {
		log.Print(err)
		user = nil
	}
	return
}

// AddFavourite for user
func (user *User) AddFavourite(recipeID string) (err error) {

	objID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		log.Print(err)
		return
	}

	newFavourite := favourite{
		RecipeID: objID,
		Added:    time.Now(),
	}
	userFilter := bson.M{
		"UserID": user.ID,
	}

	// Add the new recipe to the collection if not existing
	update := bson.M{
		"$addToSet": bson.M{
			"Favourites": newFavourite,
		},
	}

	// Try to update if the exsiting user
	_, err = database.Manager.Get(collection).UpdateOne(*database.Manager.GetContext(), userFilter, update)
	if err != nil {
		log.Print(err)
	}
	return
}
