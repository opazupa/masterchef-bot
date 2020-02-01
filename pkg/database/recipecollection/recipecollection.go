package recipecollection

import (
	"masterchef_bot/pkg/database"
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

const collection = "recipes"

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

	results := []Recipe{}

	ctx := *database.Manager.GetContext()
	cursor, err := database.Manager.Get(collection).Find(ctx, filter)
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
func GetByID(id primitive.ObjectID) *Recipe {

	filter := bson.D{
		primitive.E{
			Key: "_id", Value: id,
		},
	}
	var result Recipe

	err := database.Manager.Get(collection).FindOne(*database.Manager.GetContext(), filter).Decode(&result)
	if err != nil {
		return nil
	}

	return &result
}
