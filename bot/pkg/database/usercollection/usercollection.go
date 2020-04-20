package usercollection

import (
	"log"
	"masterchef_bot/pkg/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Favourite elem
type Favourite struct {
	RecipeID primitive.ObjectID `bson:"RecipeID"`
}

// User document
type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserName   string             `bson:"UserName"`
	TelegramID int                `bson:"TelegramID"`
	Registered time.Time          `bson:"Registered"`
	Favourites []Favourite        `bson:"Favourites"`
}

// Create new user
func Create(userName string, id int) (user *User, err error) {
	newUser := bson.M{
		"UserName":   userName,
		"TelegramId": id,
		"Registered": time.Now(),
		"Favourites": []Favourite{},
	}

	inserted, err := database.Manager.Get(database.Users).InsertOne(*database.Manager.GetContext(), newUser)
	if err != nil {
		return
	}
	user, err = getByID(inserted.InsertedID.(primitive.ObjectID))
	return
}

// GetByID from collection
func getByID(id primitive.ObjectID) (user *User, err error) {

	user = &User{}
	filter := bson.M{
		"_id": id,
	}
	err = database.Manager.Get(database.Users).FindOne(*database.Manager.GetContext(), filter).Decode(user)
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
	filter := bson.M{
		"UserName": userName,
	}

	err := database.Manager.Get(database.Users).FindOne(*database.Manager.GetContext(), filter).Decode(user)
	if err != nil {
		log.Print(err)
		user = nil
	}
	return
}

// AddFavourite for user
func (user *User) AddFavourite(recipeID string) (added bool, err error) {

	objID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		log.Print(err)
		return
	}

	userFilter := bson.M{
		"_id": user.ID,
	}

	// Add the new recipe to the collection if not existing
	update := bson.M{
		"$addToSet": bson.M{
			"Favourites": Favourite{
				RecipeID: objID,
			},
		},
	}

	// Try to update if the exsiting user
	opts, err := database.Manager.Get(database.Users).UpdateOne(*database.Manager.GetContext(), userFilter, update)
	if err != nil {
		log.Print(err)
	}
	added = opts.ModifiedCount != 0
	return
}

// RemoveFavourite for user
func (user *User) RemoveFavourite(recipeID string) (removed bool, err error) {

	objID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		log.Print(err)
		return
	}

	userFilter := bson.M{
		"_id": user.ID,
	}

	// Remove the new recipe to the collection if existing
	update := bson.M{
		"$pull": bson.M{
			"Favourites": bson.M{
				"RecipeID": objID,
			},
		},
	}

	// Try to update if the exsiting user
	opts, err := database.Manager.Get(database.Users).UpdateOne(*database.Manager.GetContext(), userFilter, update)
	if err != nil {
		log.Print(err)
	}
	removed = opts.ModifiedCount != 0
	return
}

// IsRegistered to app
func (user *User) IsRegistered() bool {
	return user != nil
}
