package inlinequery

import (
	"log"
	"masterchef_bot/pkg/googleapi"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Handle inline query for the bot
func Handle(update *tgbotapi.Update) *[]tgbotapi.InlineQueryResultArticle {

	log.Printf("[%s] Inline query: (%s)", update.InlineQuery.From.UserName, update.InlineQuery.Query)
	if update.InlineQuery.Query == "" {
		// Return maybe some trending results here
	}
	results := googleapi.SearchRecipes(update.InlineQuery.Query)
	return toInlineQueryResult(results)
}

// Convert google search results to InlineQueryResults
func toInlineQueryResult(recipes *[]googleapi.Recipe) *[]tgbotapi.InlineQueryResultArticle {

	return &[]tgbotapi.InlineQueryResultArticle{}
}
