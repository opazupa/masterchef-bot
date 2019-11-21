package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"

	"masterchef_bot/pkg/bot/command"
	"masterchef_bot/pkg/bot/inlinequery"
	"masterchef_bot/pkg/configuration"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system if present
	godotenv.Load()
}

// Main application
func main() {

	bot := configureBot()
	handleUpdates(bot)
}

// Configure bot on startup
func configureBot() *tgbotapi.BotAPI {

	configuration := configuration.Get()
	bot, err := tgbotapi.NewBotAPI(configuration.APIKey)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = configuration.DebugMode

	log.Printf("Fired up %s 🔥🔥🔥", bot.Self.UserName)
	return bot
}

// Handle received updates
func handleUpdates(bot *tgbotapi.BotAPI) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Printf("Failed to get update feed.")
	}
	for update := range updates {

		if update.InlineQuery != nil {
			results := inlinequery.Handle(&update)
			response := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				Results:       *results,
			}
			bot.AnswerInlineQuery(response)

		} else if update.Message.IsCommand() && update.Message != nil {
			msg, err := command.Handle(&update, bot.Self.UserName)
			if err == nil {
				bot.Send(msg)
			}
		}
	}
}
