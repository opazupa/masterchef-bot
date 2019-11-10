package command

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
@Mc_Recipe_Bot and then by typing recipe.

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
`,
	},
}

// Handle command for bot
func Handle(update *tgbotapi.Update) (msg *tgbotapi.MessageConfig, err error) {

	var reply tgbotapi.MessageConfig

	switch update.Message.Command() {
	case strings.ToLower(commands.Help.Key):
		reply = tgbotapi.NewMessage(update.Message.Chat.ID, commands.Help.Description)
	case strings.ToLower(commands.Start.Key):
		reply = tgbotapi.NewMessage(update.Message.Chat.ID, commands.Start.Description)
	default:
		return nil, fmt.Errorf("Unregocnized command %s from user  [%s]", update.Message.Command(), update.Message.From.UserName)
	}
	return &reply, nil
}
