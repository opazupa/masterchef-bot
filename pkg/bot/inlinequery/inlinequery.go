package inlinequery

import (
	"log"
	"masterchef_bot/pkg/duckduckgoapi"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Handle inline query for the bot
func Handle(update *tgbotapi.Update) *[]interface{} {

	log.Printf("[%s] Inline query: (%s)", update.InlineQuery.From.UserName, update.InlineQuery.Query)
	if update.InlineQuery.Query == "" {
		// Return maybe some trending results here
		return &[]interface{}{}
	}
	results := duckduckgoapi.SearchRecipes(update.InlineQuery.Query)
	return toInlineQueryResult(results)
}

// Convert google search results to InlineQueryResults
func toInlineQueryResult(recipes *[]duckduckgoapi.Recipe) *[]interface{} {

	results := make([]interface{}, len(*recipes))
	for i := 1; i < 5; i++ {
		x := tgbotapi.NewInlineQueryResultArticle(strconv.Itoa(i), "title"+string(i), "https://www.k-ruoka.fi/reseptit/palak-paneer")
		x.URL = "https://www.k-ruoka.fi/reseptit/palak-paneer"
		x.HideURL = true
		x.Description = "moi"

		results = append(results, x)
	}
	return &results
}
