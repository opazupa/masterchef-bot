package recipecollection

import (
	"fmt"
	"masterchef_bot/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Recipe document
type Recipe struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID primitive.ObjectID `bson:"UserID"`
	Name   string             `bson:"Name"`
	URL    string             `bson:"URL"`
}

const recipeCollection = "recipes"

// Add new recipe
func Add(name string, url string, userID primitive.ObjectID) (*Recipe, error) {
	newRecipe := bson.M{
		"Name":   name,
		"URL":    url,
		"UserID": userID,
	}

	db := database.Get()
	inserted, err := db.Collection(recipeCollection).InsertOne(database.GetContext(), newRecipe)
	if err != nil {
		fmt.Print("moi")
		return nil, err
	}

	return getByID(inserted.InsertedID.(primitive.ObjectID)), nil
}

// GetByUser from collection
func GetByUser(userID primitive.ObjectID) *[]Recipe {

	filter := bson.D{
		primitive.E{
			Key: "UserID", Value: userID,
		},
	}

	results := []Recipe{}

	db := database.Get()
	ctx := database.GetContext()
	cursor, err := db.Collection(recipeCollection).Find(ctx, filter)
	if err != nil {
		return &results
	}
	// Iterate through the returned cursor.
	for cursor.Next(ctx) {
		var doc Recipe
		cursor.Decode(&doc)
		results = append(results, doc)
	}

	return &results
}

// GetByID from collection
func getByID(id primitive.ObjectID) *Recipe {

	filter := bson.D{
		primitive.E{
			Key: "_id", Value: id,
		},
	}
	var result Recipe

	db := database.Get()
	err := db.Collection(recipeCollection).FindOne(database.GetContext(), filter).Decode(&result)
	if err != nil {
		return nil
	}

	return &result
}
