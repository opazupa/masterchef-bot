package command

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Handle command for bot
func Handle(update *tgbotapi.Update) {
	log.Printf("[%s] Command: %s", update.Message.From.UserName, update.Message.Command())
}
