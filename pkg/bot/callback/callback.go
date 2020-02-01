package callback

import (
	"fmt"
	"log"

	"masterchef_bot/pkg/database/recipecollection"
	selection "masterchef_bot/pkg/database/selectedrecipecollection"
	"masterchef_bot/pkg/database/usercollection"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// callbackAction type
type callbackAction struct {
	ID   string
	Text string
}

// Actions type
type Actions struct {
	RegisterAction callbackAction
	SaveAction     callbackAction
}

// RegisteredActions for callbacks
var RegisteredActions = &Actions{
	// SaveAction for save recipe buttons
	SaveAction: callbackAction{
		ID:   "1",
		Text: "Save 👊",
	},
	// RegisterAction for register user button
	RegisterAction: callbackAction{
		ID:   "2",
		Text: "Hop on 👌",
	},
}

// Handle callbackquery updates and next action
func Handle(update *tgbotapi.Update, user *usercollection.User) (msg string, nextAction *tgbotapi.EditMessageReplyMarkupConfig) {

	var replyText string

	switch update.CallbackQuery.Data {
	case RegisteredActions.SaveAction.ID:
		if user == nil {
			return "Register first to start collecting recipes.", createAction(update.CallbackQuery, nil)
		}

		// Get user's selection from database
		selectedRecipe := selection.GetByUser(user.ID)
		if selectedRecipe == nil {
			return "Something went wrong when fetching the selected recipe 🧐", createAction(update.CallbackQuery, nil)
		}

		// Save recipe to database
		_, err := recipecollection.Add(selectedRecipe.Name, selectedRecipe.URL, user.ID)

		if err == nil {
			replyText = fmt.Sprintf("Recipe saved 😛")
		} else {
			replyText = fmt.Sprintf("Failed to save the recipe ☹")
		}

	case RegisteredActions.RegisterAction.ID:
		if user != nil {
			return "You're already registered.", createAction(update.CallbackQuery, nil)
		}

		// Register user for the bot
		newUser, err := usercollection.Create(update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID)
		if err == nil {
			replyText = fmt.Sprintf("User [%s] registered 🔥", newUser.UserName)
		} else {
			replyText = fmt.Sprintf("Failed to register user [%s]", update.CallbackQuery.From.UserName)
		}

	default:
		log.Printf("Unregocnized callback (%s) from user [%s]", update.CallbackQuery.Data, update.CallbackQuery.From.UserName)
		replyText = "Unknown callback 🧐"
	}

	return replyText, createAction(update.CallbackQuery, nil)
}

// Create nee action keyboard by modifying the existing message
func createAction(callback *tgbotapi.CallbackQuery, markup *tgbotapi.InlineKeyboardMarkup) *tgbotapi.EditMessageReplyMarkupConfig {
	nextAction := tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ReplyMarkup: markup,
		},
	}
	if callback.InlineMessageID != "" {
		nextAction.InlineMessageID = callback.InlineMessageID
	}
	if callback.Message != nil {
		nextAction.ChatID = callback.Message.Chat.ID
		nextAction.MessageID = callback.Message.MessageID
	}
	return &nextAction
}
