package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"

	"masterchef_bot/pkg/bot/callback"
	"masterchef_bot/pkg/bot/command"
	"masterchef_bot/pkg/bot/inlinequery"
	"masterchef_bot/pkg/configuration"
	"masterchef_bot/pkg/database"
	"masterchef_bot/pkg/database/usercollection"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system if present
	godotenv.Load()
	database.Check()
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

		// Check if the user is registered!
		userID := getUser(update)
		registeredUser := usercollection.Get(userID)
		log.Print(registeredUser)
		log.Print(registeredUser != nil)

		if update.InlineQuery != nil {
			// When user searches recipes with inline query
			results := inlinequery.Handle(&update, registeredUser != nil)
			response := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				Results:       *results,
			}
			bot.AnswerInlineQuery(response)
			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID
			// bot.Send(msg)

		} else if update.CallbackQuery != nil {
			// When user interacts with inline buttons
			err := callback.Handle(&update)
			if err != nil {
				errMsg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Unknown callback 🧐")
				bot.Send(errMsg)
			} else {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ok")
				msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
				bot.Send(msg)
			}

		} else if update.Message.IsCommand() && update.Message != nil {
			// When user enter a command
			msg, err := command.Handle(&update, bot.Self.UserName)
			if err == nil {
				bot.Send(msg)
			}
		}
	}
}

func getUser(update tgbotapi.Update) (id *int) {

	if update.Message != nil {
		return &update.Message.From.ID
	} else if update.InlineQuery != nil {
		return &update.InlineQuery.From.ID
	} else if update.CallbackQuery != nil {
		return &update.CallbackQuery.From.ID
	} else {
		log.Print("Unable to find user id from update", update)
		return nil
	}
}
