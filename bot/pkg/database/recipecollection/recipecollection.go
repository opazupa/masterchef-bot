package recipecollection

import (
	"fmt"
	"log"
	"masterchef_bot/pkg/database"
	"masterchef_bot/pkg/database/selectedrecipecollection"
	templates "masterchef_bot/pkg/helpers"
	"time"

	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Recipe document
type Recipe struct {
	ID      primitive.ObjectID `bson:"_id"`
	UserID  primitive.ObjectID `bson:"UserID"`
	Name    string             `bson:"Name"`
	URL     string             `bson:"URL"`
	Added   time.Time          `bson:"Added"`
	Updated time.Time          `bson:"Updated, omitempty"`
}

// FavouriteRecipe extending the Recipe
type FavouriteRecipe struct {
	Recipe     `bson:",inline"`
	Favourited int `bson:"Favourited"`
}

// ToMessage from recipe with given title
func (recipe Recipe) ToMessage(header ...string) (message string) {
	return fmt.Sprintf(templates.RecipeMessage, header, recipe.Name, recipe.URL)
}

// ToMessage from favourite recipe with given title
func (recipe FavouriteRecipe) ToMessage(header ...string) (message string) {
	return fmt.Sprintf(templates.FavouriteRecipeMessage, header, recipe.Name, recipe.URL, recipe.Favourited)
}

// Add new recipe
func Add(recipe *selectedrecipecollection.SelectedRecipe) (addedRecipe *Recipe, err error) {
	newRecipe := bson.M{
		"Name":   recipe.Name,
		"URL":    recipe.URL,
		"UserID": recipe.UserID,
		"Added":  time.Now(),
	}

	inserted, err := database.Manager.Get(database.Recipes).InsertOne(*database.Manager.GetContext(), newRecipe)
	if err != nil {
		sentry.CaptureException(err)
		log.Print(err)
		return
	}
	addedRecipe, err = GetByID(inserted.InsertedID.(primitive.ObjectID))
	return
}

// GetByUser from collection
func GetByUser(userID primitive.ObjectID) (recipes *[]FavouriteRecipe) {

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"UserID": userID,
			},
		},
		{
			"$lookup": bson.M{
				"from":         database.Users,
				"localField":   "_id",
				"foreignField": "Favourites.RecipeID",
				"as":           "Favourited",
			},
		},
		{
			"$addFields": bson.M{
				"Favourited": bson.M{
					"$size": "$Favourited",
				},
			},
		},
		{
			"$sort": bson.M{
				"Name": 1,
			},
		},
	}

	ctx := *database.Manager.GetContext()
	cursor, err := database.Manager.Get(database.Recipes).Aggregate(ctx, pipeline)
	if err != nil {
		sentry.CaptureException(err)
		log.Print(err)
		return
	}

	recipes = &[]FavouriteRecipe{}
	cursor.All(ctx, recipes)
	return
}

// GetFavouritesByUser from collection
func GetFavouritesByUser(userID primitive.ObjectID) (recipes *[]FavouriteRecipe) {

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         database.Users,
				"localField":   "_id",
				"foreignField": "Favourites.RecipeID",
				"as":           "FavouritedBy",
			},
		},
		{
			"$match": bson.M{
				"FavouritedBy._id": userID,
			},
		},
		{
			"$addFields": bson.M{
				"Favourited": bson.M{
					"$size": "$FavouritedBy",
				},
			},
		},
		{
			"$sort": bson.M{
				"Name": 1,
			},
		},
	}

	ctx := *database.Manager.GetContext()
	cursor, err := database.Manager.Get(database.Recipes).Aggregate(ctx, pipeline)
	if err != nil {
		log.Print(err)
		return
	}

	recipes = &[]FavouriteRecipe{}
	cursor.All(ctx, recipes)
	return
}

// GetByID from collection
func GetByID(id primitive.ObjectID) (recipe *Recipe, err error) {

	recipe = &Recipe{}
	filter := bson.M{
		"_id": id,
	}

	err = database.Manager.Get(database.Recipes).FindOne(*database.Manager.GetContext(), filter).Decode(recipe)
	if err != nil {
		log.Print(err)
		recipe = nil
	}
	return
}

// GetRandom recipes from collection
func GetRandom(amount int) (recipes *[]FavouriteRecipe) {

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         database.Users,
				"localField":   "_id",
				"foreignField": "Favourites.RecipeID",
				"as":           "Favourited",
			},
		},
		{
			"$addFields": bson.M{
				"Favourited": bson.M{
					"$size": "$Favourited",
				},
			},
		},
		{
			"$sample": bson.M{
				"size": amount,
			},
		},
	}

	ctx := *database.Manager.GetContext()
	cursor, err := database.Manager.Get(database.Recipes).Aggregate(ctx, pipeline)
	if err != nil {
		log.Print(err)
		return
	}

	recipes = &[]FavouriteRecipe{}
	cursor.All(ctx, recipes)
	return
}

// GetMostFavourited recipes from collection
func GetMostFavourited(amount int) (recipes *[]FavouriteRecipe) {

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         database.Users,
				"localField":   "_id",
				"foreignField": "Favourites.RecipeID",
				"as":           "Favourited",
			},
		},
		{
			"$addFields": bson.M{
				"Favourited": bson.M{
					"$size": "$Favourited",
				},
			},
		},
		{
			"$sort": bson.M{
				"Favourited": -1,
			},
		},
		{"$limit": amount},
	}

	ctx := *database.Manager.GetContext()
	cursor, err := database.Manager.Get(database.Recipes).Aggregate(ctx, pipeline)
	if err != nil {
		log.Print(err)
		return
	}

	recipes = &[]FavouriteRecipe{}
	cursor.All(ctx, recipes)
	return
}
