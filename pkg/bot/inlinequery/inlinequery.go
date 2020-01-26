package inlinequery

import (
	"fmt"
	"masterchef_bot/pkg/bot/callback"
	"masterchef_bot/pkg/recipeapi"
	"strconv"
	"strings"

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

// GetRecipeInfo for an selected inlinequery recipe
func GetRecipeInfo(update *tgbotapi.Update) (name string, url string) {
	recipeResultParts := strings.Split(update.Message.Text, "\n")
	return recipeResultParts[namePosition], recipeResultParts[urlPosition]
}

// Convert recipe results to InlineQueryResults
func toInlineQueryResult(recipes *[]recipeapi.Recipe, isRegistered bool) *[]interface{} {

	titleTemplate := `
%s
%s
%s`

	results := make([]interface{}, 0)
	for i, recipe := range (*recipes)[:resultLimit] {
		results = append(results, tgbotapi.InlineQueryResultArticle{
			Type:  "article",
			ID:    strconv.Itoa(i + 1),
			Title: recipe.Title,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: fmt.Sprintf(titleTemplate, recipeResultPrefix, recipe.Title, recipe.URL),
			},
			URL:         recipe.URL,
			ThumbHeight: 8,
			ThumbWidth:  8,
			ThumbURL:    "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQcCPYO-yVEALy1NE2deQtHC2uOy091lUvRPyFWEUyE0xlgsNm8&s",
			HideURL:     true,
			Description: recipe.Description,
			ReplyMarkup: addActionButtons(isRegistered),
		})
	}

	return &results
}

func addActionButtons(isRegistered bool) *tgbotapi.InlineKeyboardMarkup {

	// Hide buttons if user is not registered
	if !isRegistered {
		return nil
	}

	saveAction := callback.RegisteredActions.SaveAction
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(saveAction.Text, fmt.Sprint(saveAction.ID)),
		),
	)
	return &keyboard
}
