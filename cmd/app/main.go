package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"

	"masterchef_bot/pkg/bot/callback"
	"masterchef_bot/pkg/bot/command"
	"masterchef_bot/pkg/bot/inlinequery"
	"masterchef_bot/pkg/configuration"
	selection "masterchef_bot/pkg/database/selectedrecipecollection"
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
		registeredUser := usercollection.GetByUserName(getUser(update))

		if update.InlineQuery != nil {
			// When user searches recipes with inline query
			results := inlinequery.Handle(&update, registeredUser != nil)
			response := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				Results:       *results,
			}
			bot.AnswerInlineQuery(response)

		} else if inlinequery.IsRecipe(&update) && registeredUser != nil {
			// When registered user selects a recipe from inline query
			recipeName, recipeURL := inlinequery.GetRecipeInfo(&update)
			selection.Save(recipeName, recipeURL, update.Message.Chat.ID, registeredUser.ID)

		} else if update.CallbackQuery != nil {
			// When user interacts with inline buttons
			replyText, nextAction := callback.Handle(&update, registeredUser)
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
			msg, err := command.Handle(&update, bot.Self.UserName, registeredUser)
			if err == nil {
				bot.Send(msg)
			}
		}
	}
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
		log.Print("Unable to find username from update", update)
		return nil
	}
}
