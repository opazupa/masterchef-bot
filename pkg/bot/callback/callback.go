package callback

import (
	"fmt"
	"log"
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
	SaveAction     callbackAction
	RegisterAction callbackAction
}

// RegisteredActions for callbacks
var RegisteredActions = &Actions{
	// SaveAction for save recipe buttons
	SaveAction: callbackAction{
		ID:   "1",
		Text: "Save 👊",
	},
	// RegisterAction for registr user button
	RegisterAction: callbackAction{
		ID:   "2",
		Text: "Hop on 👌",
	},
}

// Handle callbackquery updates
func Handle(update *tgbotapi.Update) (msg string) {

	var replyText string
	switch update.CallbackQuery.Data {
	case RegisteredActions.SaveAction.ID:
		// Save recipe to database

	case RegisteredActions.RegisterAction.ID:
		newUser, err := usercollection.Create(update.CallbackQuery.From.UserName)
		if err == nil {
			replyText = fmt.Sprintf("User [%s] registered", newUser.UserName)
		} else {
			replyText = fmt.Sprintf("Failed to register user [%s]", update.CallbackQuery.From.UserName)
		}

	default:
		log.Printf("Unregocnized callback %s from user [%s]", update.CallbackQuery.Data, update.CallbackQuery.From.UserName)
		replyText = "Unknown callback 🧐"
	}
	return replyText
}
