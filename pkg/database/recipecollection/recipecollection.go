package recipecollection

import (
	"fmt"
	"log"
	"masterchef_bot/pkg/database"
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
func Add(name string, url string, userID primitive.ObjectID) (*Recipe, error) {
	newRecipe := bson.M{
		"Name":   name,
		"URL":    url,
		"UserID": userID,
		"Added":  time.Now(),
	}

	inserted, err := database.Manager.Get(collection).InsertOne(*database.Manager.GetContext(), newRecipe)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return GetByID(inserted.InsertedID.(primitive.ObjectID)), nil
}

// GetByUser from collection
func GetByUser(userID primitive.ObjectID) *[]Recipe {

	filter := bson.D{
		primitive.E{
			Key: "UserID", Value: userID,
		},
	}

	var results []Recipe

	ctx := *database.Manager.GetContext()
	cursor, err := database.Manager.Get(collection).Find(ctx, filter)
	if err != nil {
		log.Print(err)
		return &results
	}
	cursor.All(ctx, &results)
	return &results
}

// GetByID from collection
func GetByID(id primitive.ObjectID) *Recipe {

	filter := bson.D{
		primitive.E{
			Key: "_id", Value: id,
		},
	}
	var result Recipe

	err := database.Manager.Get(collection).FindOne(*database.Manager.GetContext(), filter).Decode(&result)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &result
}

// GetRandom recipes from collection
func GetRandom(amount int) *[]Recipe {

	pipeline := []bson.M{{
		"$sample": bson.D{{Key: "size", Value: amount}},
	}}

	var results []Recipe

	ctx := *database.Manager.GetContext()
	cursor, err := database.Manager.Get(collection).Aggregate(ctx, pipeline)
	if err != nil {
		log.Print(err)
		return &results
	}
	cursor.All(ctx, &results)
	return &results
}
