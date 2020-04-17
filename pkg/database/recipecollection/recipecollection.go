package recipecollection

import (
	"fmt"
	"log"
	"masterchef_bot/pkg/database"
	"masterchef_bot/pkg/database/selectedrecipecollection"
	templates "masterchef_bot/pkg/helpers"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Recipe document
type Recipe struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID primitive.ObjectID `bson:"UserID"`
	Name   string             `bson:"Name"`
	URL    string             `bson:"URL"`
	Added  time.Time          `bson:"Added"`
}

const (
	collection = "recipes"
)

// ToMessage from recipe with given title
func (recipe *Recipe) ToMessage(header string) (message string) {
	return fmt.Sprintf(templates.RecipeMessage, header, recipe.Name, recipe.URL)
}

// Add new recipe
func Add(recipe *selectedrecipecollection.SelectedRecipe) (addedRecipe *Recipe, err error) {
	newRecipe := bson.M{
		"Name":   recipe.Name,
		"URL":    recipe.URL,
		"UserID": recipe.UserID,
		"Added":  time.Now(),
	}

	inserted, err := database.Manager.Get(collection).InsertOne(*database.Manager.GetContext(), newRecipe)
	if err != nil {
		log.Print(err)
		return
	}
	addedRecipe, err = GetByID(inserted.InsertedID.(primitive.ObjectID))
	return
}

// GetByUser from collection
func GetByUser(userID primitive.ObjectID) (recipes *[]Recipe) {

	filter := bson.M{
		"UserID": userID,
	}

	ctx := *database.Manager.GetContext()
	cursor, err := database.Manager.Get(collection).Find(ctx, filter)
	if err != nil {
		log.Print(err)
		return
	}

	recipes = &[]Recipe{}
	cursor.All(ctx, recipes)
	return
}

// GetByID from collection
func GetByID(id primitive.ObjectID) (recipe *Recipe, err error) {

	recipe = &Recipe{}
	filter := bson.M{
		"_id": id,
	}

	err = database.Manager.Get(collection).FindOne(*database.Manager.GetContext(), filter).Decode(recipe)
	if err != nil {
		log.Print(err)
		recipe = nil
	}
	return
}

// GetRandom recipes from collection
func GetRandom(amount int) (recipes *[]Recipe) {

	pipeline := []bson.M{{
		"$sample": bson.M{
			"size": amount,
		},
	}}

	ctx := *database.Manager.GetContext()
	cursor, err := database.Manager.Get(collection).Aggregate(ctx, pipeline)
	if err != nil {
		log.Print(err)
		return
	}

	recipes = &[]Recipe{}
	cursor.All(ctx, recipes)
	return
}
