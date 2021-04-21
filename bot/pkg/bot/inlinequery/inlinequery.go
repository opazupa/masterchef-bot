package inlinequery

import (
	"masterchef_bot/pkg/bot/callback"
	"masterchef_bot/pkg/recipeapi"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"

	selection "masterchef_bot/pkg/database/selectedrecipecollection"
	"masterchef_bot/pkg/database/usercollection"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	prefixPosition     = 0
	namePosition       = 1
	urlPosition        = 2
	resultLimit        = 10
	recipeResultPrefix = "🍕🍕 Here you go:"
)

// Handle inline query for the bot
func Handle(update *tgbotapi.Update, isRegistered bool) *[]interface{} {

	// Ignore empty queries
	if update.InlineQuery.Query == "" {
		return &[]interface{}{}
	}
	results := recipeapi.SearchRecipes(update.InlineQuery.Query)
	return toInlineQueryResult(results, isRegistered)
}

// SaveSelectedRecipe from inlinequery
func SaveSelectedRecipe(update *tgbotapi.Update, user *usercollection.User) {

	recipeName, recipeURL := getRecipeInfo(update)
	selection.Save(recipeName, recipeURL, update.Message.Chat.ID, user.ID)
}

// IsRecipe an selected inlinequery recipe
func IsRecipe(update *tgbotapi.Update) bool {

	// Inlinequery result has a specific prefix on message
	if update.Message != nil {
		recipeResultParts := strings.Split(update.Message.Text, "\n")
		if recipeResultParts[prefixPosition] == recipeResultPrefix {
			return true
		}
	}
	return false
}

// getRecipeInfo for an selected inlinequery recipe
func getRecipeInfo(update *tgbotapi.Update) (name string, url string) {
	recipeResultParts := strings.Split(update.Message.Text, "\n")
	return strings.TrimSpace(recipeResultParts[namePosition]), strings.TrimSpace(recipeResultParts[urlPosition])
}

// Convert recipe results to InlineQueryResults
func toInlineQueryResult(recipes *[]recipeapi.Recipe, isRegistered bool) *[]interface{} {

	resultAmount := funk.MinInt([]int{len(*recipes), resultLimit}).(int)
	results := make([]interface{}, 0)
	for i, recipe := range (*recipes)[:resultAmount] {
		results = append(results, tgbotapi.InlineQueryResultArticle{
			Type:  "article",
			ID:    strconv.Itoa(i + 1),
			Title: recipe.Title,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: recipe.ToMessage(recipeResultPrefix),
			},
			URL:         recipe.URL,
			ThumbHeight: 8,
			ThumbWidth:  8,
			ThumbURL:    "https://cmkt-image-prd.freetls.fastly.net/0.1.0/ps/7519111/600/400/m2/fpnw/wm0/chef-hat-illustration-for-cooking-logo-with-love-element-.jpg?1577620994&s=54f7c96e07ef7b7479f9606910bc167c",
			HideURL:     true,
			Description: recipe.Description,
			ReplyMarkup: addSaveButton(isRegistered),
		})
	}

	return &results
}

// Add save button if user is registered
func addSaveButton(isRegistered bool) *tgbotapi.InlineKeyboardMarkup {

	// Hide buttons if user is not registered
	if !isRegistered {
		return nil
	}
	return callback.AddActions([]int{callback.SaveAction})
}
