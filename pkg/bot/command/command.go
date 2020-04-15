package command

import (
	"fmt"
	"strings"

	"masterchef_bot/pkg/bot/callback"
	"masterchef_bot/pkg/database/recipecollection"
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
	Help   command
	Random command
	Start  command
}

// Configured commands
var commands = &list{
	Help: command{
		Key: "help",
		Description: `
How can I help you *Sir?*

*Start* with /start command.

*Random recipe* with /random command.

*Recipe search*
Search for recipes by calling
''@%s'' and then by typing recipe.

`,
	},
	Random: command{
		Key:         "random",
		Description: `*Here's the random recipe for you!* 👌`,
	},
	Start: command{
		Key: "start",
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

	case commands.Help.Key:
		reply = tgbotapi.NewMessage(update.Message.Chat.ID, strings.ReplaceAll(fmt.Sprintf(commands.Help.Description, botName), "''", "`"))

	case commands.Random.Key:
		// Get a random recipe
		if recipes := *recipecollection.GetRandom(1); len(recipes) > 0 {
			reply = tgbotapi.NewMessage(update.Message.Chat.ID, recipes[0].ToMessage(commands.Random.Description))
		}
	case commands.Start.Key:
		reply = tgbotapi.NewMessage(update.Message.Chat.ID, strings.ReplaceAll(fmt.Sprintf(commands.Start.Description, botName), "''", "`"))

		// Give option to register to new users
		if user == nil {
			reply.ReplyMarkup = addRegisterButton()
		}
	default:
		return nil, fmt.Errorf("Unregocnized command %s from user [%s]", update.Message.Command(), update.Message.From.UserName)
	}
	reply.ParseMode = "Markdown"
	return &reply, nil
}

// Adds button for user registration
func addRegisterButton() *tgbotapi.InlineKeyboardMarkup {

	registerAction := callback.RegisteredActions.RegisterAction
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(registerAction.Text, registerAction.ID),
		),
	)
	return &keyboard
}
