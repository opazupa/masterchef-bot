package main

import (
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/thoas/go-funk"

	"masterchef_bot/pkg/bot/callback"
	"masterchef_bot/pkg/bot/command"
	"masterchef_bot/pkg/bot/inlinequery"
	"masterchef_bot/pkg/configuration"
	"masterchef_bot/pkg/database/usercollection"
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

	// Configure sentry
	configureSentry(configuration)

	// Configure bot
	bot, err := tgbotapi.NewBotAPI(configuration.APIKey)
	if err != nil {
		sentry.CaptureException(err)
		log.Panic(err)
	}
	bot.Debug = configuration.DebugMode

	log.Printf("Fired up %s 🔥🔥🔥", bot.Self.UserName)
	return bot
}

// Configure sentry.io 🔥😈
func configureSentry(configuration *configuration.Configuration) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         configuration.SentryDsn,
		Debug:       configuration.DebugMode,
		Environment: configuration.SentryEnv,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

// Handle received updates
func handleUpdates(bot *tgbotapi.BotAPI) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Failed to get update feed.")
	}

	for update := range updates {

		// Check if the user is registered!
		user := usercollection.GetByUserName(getUser(update))

		if update.InlineQuery != nil {
			// When user searches recipes with inline query
			results := inlinequery.Handle(&update, user.IsRegistered())
			response := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				Results:       *results,
			}
			bot.AnswerInlineQuery(response)

		} else if inlinequery.IsRecipe(&update) && user.IsRegistered() {
			// When registered user selects a recipe from inline query
			inlinequery.SaveSelectedRecipe(&update, user)

		} else if update.CallbackQuery != nil {
			// When user interacts with inline buttons
			replyText, nextAction := callback.Handle(&update, user)
			response := tgbotapi.CallbackConfig{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            replyText,
			}

			// Update inline action buttons
			if nextAction != nil {
				bot.Send(nextAction)
			}
			bot.AnswerCallbackQuery(response)

		} else if update.Message != nil && update.Message.IsCommand() {
			// When user enters a command
			messages, err := command.Handle(&update, bot.Self.UserName, user)
			if err == nil {
				funk.ForEach(*messages, func(message tgbotapi.MessageConfig) { bot.Send(message) })
			}
		}
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
}

// Get username from the update
func getUser(update tgbotapi.Update) (id *string) {

	if update.Message != nil {
		return &update.Message.From.UserName
	} else if update.InlineQuery != nil {
		return &update.InlineQuery.From.UserName
	} else if update.CallbackQuery != nil {
		return &update.CallbackQuery.From.UserName
	} else if update.EditedMessage != nil {
		return &update.EditedMessage.From.UserName
	} else {
		sentry.CaptureMessage(fmt.Sprint("Unable to find username from update", update))
		log.Print("Unable to find username from update", update)
		return nil
	}
}
