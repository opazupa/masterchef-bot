package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"

	"masterchef_bot/internal/configuration"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system if present
	godotenv.Load()
}

func main() {

	configuration := configuration.Get()
	bot, err := tgbotapi.NewBotAPI(configuration.APIKey)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = configuration.DebugMode

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
