package inlinequery

import (
	"fmt"
	"log"
	"masterchef_bot/pkg/bot/callback"
	"masterchef_bot/pkg/recipeapi"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	resultLimit = 10
)

// Handle inline query for the bot
func Handle(update *tgbotapi.Update, isRegistered bool) *[]interface{} {

	log.Printf("[%s] Inline query: (%s)", update.InlineQuery.From.UserName, update.InlineQuery.Query)
	if update.InlineQuery.Query == "" {
		// Return maybe some trending results here
		return &[]interface{}{}
	}
	results := recipeapi.SearchRecipes(update.InlineQuery.Query)
	return toInlineQueryResult(results, isRegistered)
}

// Convert recipe results to InlineQueryResults
func toInlineQueryResult(recipes *[]recipeapi.Recipe, isRegistered bool) *[]interface{} {

	titleTemplate := `
%s
%s`

	results := make([]interface{}, 0)
	for i, recipe := range (*recipes)[:resultLimit] {
		results = append(results, tgbotapi.InlineQueryResultArticle{
			Type:  "article",
			ID:    strconv.Itoa(i + 1),
			Title: recipe.Title,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: fmt.Sprintf(titleTemplate, recipe.Title, recipe.URL),
			},
			URL:         recipe.URL,
			ThumbHeight: 8,
			ThumbWidth:  8,
			ThumbURL:    "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQcCPYO-yVEALy1NE2deQtHC2uOy091lUvRPyFWEUyE0xlgsNm8&s",
			HideURL:     true,
			Description: recipe.Description,
			ReplyMarkup: addActionButtons(isRegistered, recipe),
		})
	}

	return &results
}

func addActionButtons(isRegistered bool, recipe recipeapi.Recipe) *tgbotapi.InlineKeyboardMarkup {

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
