package callback

import (
	"fmt"
	"log"
	"strings"

	"masterchef_bot/pkg/database/recipecollection"
	selection "masterchef_bot/pkg/database/selectedrecipecollection"
	"masterchef_bot/pkg/database/usercollection"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// callbackAction type
type callbackAction struct {
	ID   string
	Text string
}

// Actions type
type Actions struct {
	RegisterAction  callbackAction
	SaveAction      callbackAction
	FavouriteAction callbackAction
}

// RegisteredActions for callbacks
var RegisteredActions = &Actions{
	// SaveAction for save recipe buttons
	SaveAction: callbackAction{
		ID:   "1",
		Text: "Save 😛",
	},
	// RegisterAction for register user button
	RegisterAction: callbackAction{
		ID:   "2",
		Text: "Hop on 👌",
	},
	// FavouriteAction for collecting fav recipes
	FavouriteAction: callbackAction{
		ID:   "3",
		Text: "Add to favourites 👍",
	},
}

const (
	// ActionDelimeter for action mappings
	actionDelimeter  string = ","
	actionIDPosition        = 0
	otherIDsPosotion        = 1
)

// CreateButton for inlinekeyboard
func (action callbackAction) CreateButton(otherIds ...string) *tgbotapi.InlineKeyboardMarkup {

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				action.Text,
				fmt.Sprint(
					action.ID,
					actionDelimeter,
					strings.Join(otherIds, actionDelimeter),
				)),
		),
	)
	return &keyboard
}

// Handle callbackquery updates and next action
func Handle(update *tgbotapi.Update, user *usercollection.User) (replyMessage string, nextAction *tgbotapi.EditMessageReplyMarkupConfig) {

	// Default to shot the existing one
	nextAction = nil

	action, targetIDs := getActionInfo(update.CallbackQuery)
	switch action {

	/*
		Save Action.
		-------------
		Next Action: Favourite Action
	*/
	case RegisteredActions.SaveAction.ID:
		if user == nil {
			replyMessage = "Register first to start collecting recipes."
			break
		}

		// Get user's selection from database
		selectedRecipe := selection.GetByUser(user.ID)
		if selectedRecipe == nil {
			replyMessage = "Something went wrong when fetching the selected recipe 🧐"
			break
		}

		// Save recipe to database
		_, err := recipecollection.Add(selectedRecipe)

		if err == nil {
			replyMessage = fmt.Sprintf("Recipe saved 😛")
		} else {
			replyMessage = fmt.Sprintf("Failed to save the recipe 😕")
		}

	/*
		Register Action.
		-------------
		Next Action: empty
	*/
	case RegisteredActions.RegisterAction.ID:

		// With empty content
		nextAction = createNextAction(update.CallbackQuery, nil)
		if user != nil {
			replyMessage = "You're already registered."
			break
		}

		// Register user for the bot
		newUser, err := usercollection.Create(update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID)
		if err == nil {
			replyMessage = fmt.Sprintf("User [%s] registered 🔥", newUser.UserName)
		} else {
			replyMessage = fmt.Sprintf("Failed to register user [%s]", update.CallbackQuery.From.UserName)
		}

	/*
		Save Action.
		-------------
		Next Action: nil
	*/
	case RegisteredActions.FavouriteAction.ID:
		// TODO olli how to set it as favourite
		if user == nil {
			replyMessage = "Register first to start adding favourites."
			break
		}
		if recipeID, ok := funk.Head(targetIDs).(primitive.ObjectID); ok {

			err := user.AddFavourite(recipeID)
			if err != nil || !ok {
				replyMessage = "Something went wrong when saving favourite recipe 🧐"
				break
			}

			replyMessage = "Recipe favourited 💟"
		}

	default:
		log.Printf("Unregocnized callback (%s) from user [%s]", update.CallbackQuery.Data, update.CallbackQuery.From.UserName)
		replyMessage = "Unknown callback 🧐"
	}
	return
}

// Create next action keyboard by modifying the existing message
func createNextAction(callback *tgbotapi.CallbackQuery, content *tgbotapi.InlineKeyboardMarkup) *tgbotapi.EditMessageReplyMarkupConfig {
	nextAction := tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ReplyMarkup: content,
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

// Get action info from the callbackQuery
func getActionInfo(callbackQuery *tgbotapi.CallbackQuery) (actionID string, otherIDs []string) {
	// Actions template is <ActionId>,<otherId1>,<otherId2>
	var actionParts = strings.Split(callbackQuery.Data, actionDelimeter)
	return actionParts[actionIDPosition], actionParts[otherIDsPosotion:]
}
