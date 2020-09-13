package callback

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"masterchef_bot/pkg/database/recipecollection"
	selection "masterchef_bot/pkg/database/selectedrecipecollection"
	"masterchef_bot/pkg/database/usercollection"

	"github.com/getsentry/sentry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/thoas/go-funk"
)

// Action type
type action struct {
	ID             int
	Text           string
	NextActions    []int
	PersistOnClick bool
}

// Actions
const (
	RegisterAction    int = 1
	SaveAction        int = 2
	FavouriteAction   int = 3
	UnfavouriteAction int = 4
)

var registeredActions map[int]action = map[int]action{
	// Register Action for register user button
	RegisterAction: {
		ID:             RegisterAction,
		Text:           "Hop on 👌",
		NextActions:    nil,
		PersistOnClick: false,
	},
	// Save Action for save recipe buttons
	SaveAction: {
		ID:             SaveAction,
		Text:           "Save 😛",
		NextActions:    []int{FavouriteAction, UnfavouriteAction},
		PersistOnClick: false,
	},
	// Favourite Action for collecting fav recipes
	FavouriteAction: {
		ID:             FavouriteAction,
		Text:           "Favourite 👍",
		PersistOnClick: true,
	},
	// Unfavourite Action for cleaning fav recipes
	UnfavouriteAction: {
		ID:             UnfavouriteAction,
		Text:           "Unfavourite ❌",
		PersistOnClick: true,
	},
}

const (
	// ActionDelimeter for action mappings
	actionDelimeter  = ","
	actionIDPosition = 0
	otherIDsPosition = 1
)

// AddActions to inlinekeyboard
func AddActions(actionIds []int, otherIds ...string) *tgbotapi.InlineKeyboardMarkup {

	if actionIds == nil || !funk.Any(actionIds) {
		return nil
	}

	var buttons = funk.Map(actionIds, func(id int) tgbotapi.InlineKeyboardButton {
		action := registeredActions[id]
		return tgbotapi.NewInlineKeyboardButtonData(
			action.Text,
			fmt.Sprint(
				action.ID,
				actionDelimeter,
				strings.Join(otherIds, actionDelimeter),
			))
	}).([]tgbotapi.InlineKeyboardButton)

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
	return &keyboard
}

// Handle callbackquery updates and next action
func Handle(update *tgbotapi.Update, user *usercollection.User) (replyMessage string, nextAction *tgbotapi.EditMessageReplyMarkupConfig) {

	// Default to clear the actions
	action, targetIDs := getActionInfo(update.CallbackQuery)

	switch action {

	/*
		Register Action.
		-------------
		Next Action: empty
	*/
	case RegisterAction:

		if user.IsRegistered() {
			replyMessage = "You're already registered."
			break
		}

		// Register user for the bot
		newUser, err := usercollection.Create(update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID)
		if err == nil {
			replyMessage = fmt.Sprintf("User [%s] registered 🔥", newUser.UserName)
		} else {
			sentry.CaptureException(err)
			replyMessage = fmt.Sprintf("Failed to register user [%s]", update.CallbackQuery.From.UserName)
		}

	/*
		Save Action.
		-------------
		Next Action: Favourite Action, Unfavourite Action
	*/
	case SaveAction:
		if !user.IsRegistered() {
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
		recipe, err := recipecollection.Add(selectedRecipe)

		if err == nil {
			replyMessage = fmt.Sprintf("Recipe saved 😛")
		} else {
			sentry.CaptureException(err)
			replyMessage = fmt.Sprintf("Failed to save the recipe 😕")
			break
		}
		nextAction = getNextAction(action, err, update.CallbackQuery, recipe.ID.Hex())

	/*
		Favourite Action.
		-------------
		Next Action: no change
	*/
	case FavouriteAction:

		nextAction = getNextAction(action, nil, update.CallbackQuery)

		if !user.IsRegistered() {
			replyMessage = "Register first to start adding favourites."
			break
		}

		added, err := user.AddFavourite(funk.Head(targetIDs).(string))
		if err != nil {
			sentry.CaptureException(err)
			replyMessage = "Something went wrong when saving favourite recipe 🧐"
			break
		}
		if added {
			replyMessage = "Recipe favourited 💟"
		} else {
			replyMessage = "Recipe already favourited!"
		}

	/*
		Unfavourite Action.
		-------------
		Next Action: no cahnge
	*/
	case UnfavouriteAction:

		nextAction = getNextAction(action, nil, update.CallbackQuery)

		if !user.IsRegistered() {
			replyMessage = "Register first to start adding favourites."
			break
		}

		_, err := user.RemoveFavourite(funk.Head(targetIDs).(string))
		if err != nil {
			sentry.CaptureException(err)
			replyMessage = "Something went wrong when removing favourite recipe 🧐"
			break
		}
		replyMessage = "Recipe unfavourited"

	default:
		log.Printf("Unregocnized callback (%s) from user [%s]", update.CallbackQuery.Data, update.CallbackQuery.From.UserName)
		replyMessage = "Unknown callback 🧐"
	}
	return
}

type nextActionType int

// Create next action keyboard by modifying the existing message
func getNextAction(actionID int, err error, callback *tgbotapi.CallbackQuery, otherIds ...string) *tgbotapi.EditMessageReplyMarkupConfig {

	action := registeredActions[actionID]
	// Return no changes -> action remains visible
	if action.PersistOnClick || err != nil {
		return nil
	}

	// Set the properties
	nextAction := tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ReplyMarkup: AddActions(action.NextActions, otherIds...),
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
func getActionInfo(callbackQuery *tgbotapi.CallbackQuery) (actionID int, otherIDs []string) {
	// Actions template is <ActionId>,<otherId1>,<otherId2>
	var actionParts = strings.Split(callbackQuery.Data, actionDelimeter)
	otherIDs = actionParts[otherIDsPosition:]
	actionID, err := strconv.Atoi(actionParts[actionIDPosition])

	if err != nil {
		sentry.CaptureException(err)
		log.Println(err)
	}
	return
}
