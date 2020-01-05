package callback

import (
	"fmt"

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
func Handle(update *tgbotapi.Update) (err error) {

	fmt.Print(*update.CallbackQuery)
	switch update.CallbackQuery.Data {
	case RegisteredActions.SaveAction.ID:

	case RegisteredActions.RegisterAction.ID:

	default:
		return fmt.Errorf("Unregocnized callback  %s from user [%s]", update.CallbackQuery.Data, update.CallbackQuery.From.UserName)
	}
	return nil
}
