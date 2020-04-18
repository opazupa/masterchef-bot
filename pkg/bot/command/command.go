package command

import (
	"fmt"
	"log"
	"strings"

	"masterchef_bot/pkg/bot/callback"
	"masterchef_bot/pkg/database/recipecollection"
	"masterchef_bot/pkg/database/usercollection"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/thoas/go-funk"
)

// Command type
type command struct {
	Key         string
	Description string
	NextActions []int
}

// commands
const (
	helpCommand   = "help"
	randomCommand = "random"
	startCommand  = "start"
	top3Command   = "top3"
)

var registeredCommands map[string]command = map[string]command{
	// Help command
	helpCommand: {
		Key: helpCommand,
		Description: `
		How can I help you *Sir?*
		
		*Start* with /start command.
		
		*Random recipe* with /random command.
		
		*Top3 recipes* with /top3 command.
		
		*Recipe search*
		Search for recipes by calling
		''@%s'' and then by typing recipe.
		`,
		NextActions: nil,
	},
	// Random command
	randomCommand: {
		Key:         randomCommand,
		NextActions: []int{callback.FavouriteAction, callback.UnfavouriteAction},
	},
	// Start command
	startCommand: {
		Key: startCommand,
		Description: `
		*Hi*!
		
		I’m the *Masterchef* bot on your service!👌
		
		Register and start building
		your own masterchef recipe book. 👇
		*Let’s start cooking ay?* 🔥
		
		''''''
		By doing that you accept @%s
		to store your name and telegram id.
		''''''
		`,
		NextActions: []int{callback.RegisterAction},
	},
	// Top3 command
	top3Command: {
		Key:         top3Command,
		NextActions: []int{callback.FavouriteAction, callback.UnfavouriteAction},
	},
}

// Handle command for bot
func Handle(update *tgbotapi.Update, botName string, user *usercollection.User) (reply *[]tgbotapi.MessageConfig, err error) {

	var messages []tgbotapi.MessageConfig
	command, found := registeredCommands[update.Message.Command()]

	if !found {
		err = fmt.Errorf("Unregocnized command %s from user [%s]", update.Message.Command(), update.Message.From.UserName)
		log.Print(err)
	}

	switch command.Key {

	/*
		Help command
		-------------
		Next Action: nil
	*/
	case helpCommand:
		messages = append(messages, command.messageFromDescription(update.Message.Chat.ID, botName))

	/*
		Random command
		-------------
		Next Action: Favourite action, Unfavourite action
	*/
	case randomCommand:
		// Get a random recipe
		if recipes := *recipecollection.GetRandom(1); funk.Any(recipes) {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, funk.Head(recipes).(recipecollection.Recipe).ToMessage())

			// Add action
			message.ReplyMarkup = command.getNextAction(funk.Head(recipes).(recipecollection.Recipe).ID.Hex())
			messages = append(messages, message)
		} else {
			err = fmt.Errorf("No recipes returned for the random one")
		}

	/*
		Start command
		-------------
		Next Action: Register action
	*/
	case startCommand:
		message := command.messageFromDescription(update.Message.Chat.ID, botName)

		// Give option to register to new users
		if user == nil {
			message.ReplyMarkup = command.getNextAction()
		}

		messages = append(messages, message)

	/*
		Top3 command
		-------------
		Next Action: Favourite action, Unfavourite action
	*/
	case top3Command:
		// TODO Olli get most popular recipes
		topRecipes := recipecollection.GetMostFavourited(3)
		funk.ForEach(*topRecipes, func(topRecipe recipecollection.Recipe) {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, topRecipe.ToMessage())
			message.ReplyMarkup = command.getNextAction(topRecipe.ID.Hex())
			messages = append(messages, message)
		})
	}

	// Set markdown rendering
	messages = funk.Map(messages, func(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
		message.ParseMode = "Markdown"
		return message
	}).([]tgbotapi.MessageConfig)
	return &messages, err
}

// Create message from command description
func (command command) messageFromDescription(chatID int64, params ...interface{}) (message tgbotapi.MessageConfig) {
	message = tgbotapi.NewMessage(chatID, strings.ReplaceAll(fmt.Sprintf(command.Description, params...), "''", "`"))
	return
}

// Get next action for the command
func (command command) getNextAction(otherIds ...string) *tgbotapi.InlineKeyboardMarkup {
	return callback.AddActions(command.NextActions, otherIds...)
}
