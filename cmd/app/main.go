package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"masterchef_bot/pkg/bot/command"
	"masterchef_bot/pkg/bot/inlinequery"
	"masterchef_bot/pkg/configuration"
	"masterchef_bot/pkg/database"
	"masterchef_bot/pkg/database/user"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system if present
	godotenv.Load()
}

// Main application
func main() {

	bot := configureBot()
	db := database.Get()
	handleUpdates(bot, db)
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
func handleUpdates(bot *tgbotapi.BotAPI, db *mongo.Database) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("Failed to get update feed.")
	}

	for update := range updates {

		// Check if the user is registered!
		registeredUser := user.Get(db, update.Message.From.ID)

		if update.InlineQuery != nil {
			// When user searches recipes with inline query
			results := inlinequery.Handle(&update, registeredUser != nil)
			response := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				Results:       *results,
			}
			bot.AnswerInlineQuery(response)

		} else if update.CallbackQuery != nil {
			// When user interacts with inline buttons
			log.Print("Call back query:", update.CallbackQuery.Data)

		} else if update.Message.IsCommand() && update.Message != nil {
			// When user enter a command
			msg, err := command.Handle(&update, db, bot.Self.UserName)
			if err == nil {
				bot.Send(msg)
			}
		}
	}
}
