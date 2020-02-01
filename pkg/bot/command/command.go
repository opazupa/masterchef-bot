package command

import (
	"fmt"
	"strings"

	"masterchef_bot/pkg/bot/callback"
	"masterchef_bot/pkg/database/usercollection"

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
How can I help you *Sir?*

*Start* with /start command.

*Recipe search*
Search for recipes by calling
''@%s'' and then by typing recipe.

`,
	},
	Start: command{
		Key: "Start",
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
	},
}

// Handle command for bot
func Handle(update *tgbotapi.Update, botName string, user *usercollection.User) (msg *tgbotapi.MessageConfig, err error) {

	var reply tgbotapi.MessageConfig

	switch update.Message.Command() {
	case strings.ToLower(commands.Help.Key):
		reply = tgbotapi.NewMessage(update.Message.Chat.ID, strings.ReplaceAll(fmt.Sprintf(commands.Help.Description, botName), "''", "`"))
	case strings.ToLower(commands.Start.Key):
		reply = tgbotapi.NewMessage(update.Message.Chat.ID, strings.ReplaceAll(fmt.Sprintf(commands.Start.Description, botName), "''", "`"))

		// Give option to register to new users
		if user == nil {
			reply.ReplyMarkup = addActionButtons()
		}
	default:
		return nil, fmt.Errorf("Unregocnized command %s from user [%s]", update.Message.Command(), update.Message.From.UserName)
	}
	reply.ParseMode = "Markdown"
	return &reply, nil
}

func addActionButtons() *tgbotapi.InlineKeyboardMarkup {

	registerAction := callback.RegisteredActions.RegisterAction
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(registerAction.Text, registerAction.ID),
		),
	)
	return &keyboard
}
