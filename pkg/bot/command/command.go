package command

import (
	"fmt"
	"strings"

	"masterchef_bot/pkg/bot/actionbuttons"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
)

// Command type
type command struct {
	Key         string
	Description string
}

// List of Commands
type list struct {
	Help  command
	Start command
}

// Configured commands
var commands = &list{
	Help: command{
		Key: "Help",
		Description: `
How can I help you Sir?


Start with /start command.

Recipe search
Search for recipes by calling
@%s and then by typing recipe.

Services
- list of service coming up here
`,
	},
	Start: command{
		Key: "Start",
		Description: `
Hi!

I’m the Masterchef bot on your service!👌
Let’s start cooking ay? 🔥

By registering you can start building
your own masterchef recipe book. 👇

By doing that you accept @Mc_Recipe_Bot
to store your name and telegram id.
`,
	},
}

// Handle command for bot
func Handle(update *tgbotapi.Update, db *mongo.Database, botName string) (msg *tgbotapi.MessageConfig, err error) {

	var reply tgbotapi.MessageConfig

	switch update.Message.Command() {
	case strings.ToLower(commands.Help.Key):
		reply = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(commands.Help.Description, botName))
	case strings.ToLower(commands.Start.Key):
		reply = tgbotapi.NewMessage(update.Message.Chat.ID, commands.Start.Description)
		// Give option to register to new users
		if true {
			reply.ReplyMarkup = addActionButtons()
		}
	default:
		return nil, fmt.Errorf("Unregocnized command %s from user [%s]", update.Message.Command(), update.Message.From.UserName)
	}
	return &reply, nil
}

func addActionButtons() *tgbotapi.InlineKeyboardMarkup {
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(actionbuttons.RegisterAction, ""),
		),
	)
	return &keyboard
}
