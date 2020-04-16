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
func Handle(update *tgbotapi.Update, botName string, user *usercollection.User) (reply *tgbotapi.MessageConfig, err error) {

	var message = tgbotapi.MessageConfig{
		ParseMode: "Markdown",
	}

	switch update.Message.Command() {

	/*
		Help command
		-------------
		Next Action: nil
	*/
	case commands.Help.Key:
		message = tgbotapi.NewMessage(update.Message.Chat.ID, finalizedMarkdown(commands.Help.Description, botName))

	/*
		Random command
		-------------
		Next Action: Favourite action
	*/
	case commands.Random.Key:
		// Get a random recipe
		var a *[]int = &[]int{1, 2, 3}
		log.Print(a)
		log.Print(funk.Any(a))
		log.Print(funk.Any(&a))
		if recipes := *recipecollection.GetRandom(1); true {
			log.Print(recipes)
			log.Print(funk.Any((&recipes)))
			message = tgbotapi.NewMessage(update.Message.Chat.ID, (recipes)[0].ToMessage(commands.Random.Description))
		} else {
			err = fmt.Errorf("No recipes returned for the random one")
		}

	/*
		Start command
		-------------
		Next Action: Register action
	*/
	case commands.Start.Key:
		message = tgbotapi.NewMessage(update.Message.Chat.ID, finalizedMarkdown(commands.Start.Description, botName))

		// Give option to register to new users
		if user == nil {
			message.ReplyMarkup = callback.RegisteredActions.RegisterAction.CreateButton()
		}
	default:
		err = fmt.Errorf("Unregocnized command %s from user [%s]", update.Message.Command(), update.Message.From.UserName)
	}

	return &message, err
}

// Finalize markdown with correct chars
func finalizedMarkdown(markdown string, params ...interface{}) string {
	return strings.ReplaceAll(fmt.Sprintf(markdown, params...), "''", "`")
}
