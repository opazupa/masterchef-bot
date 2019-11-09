package inlinequery

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Handle inline query for the bot
func Handle(update *tgbotapi.Update) {
	log.Printf("[%s] Inline query: (%s)", update.InlineQuery.From.UserName, update.InlineQuery.Query)
}
