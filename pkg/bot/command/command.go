package command

import (
	"fmt"
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
	Top3   command
}

// Configured commands
var commands = &list{
	Help: command{
		Key: "help",
		Description: `
How can I help you *Sir?*

*Start* with /start command.

*Random recipe* with /random command.

*Top3 recipes* with /top3 command.

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
	Top3: command{
		Key:         "top3",
		Description: "",
	},
}

// Handle command for bot
func Handle(update *tgbotapi.Update, botName string, user *usercollection.User) (reply *[]tgbotapi.MessageConfig, err error) {

	var messages []tgbotapi.MessageConfig

	switch update.Message.Command() {

	/*
		Help command
		-------------
		Next Action: nil
	*/
	case commands.Help.Key:
		messages = append(messages, tgbotapi.NewMessage(update.Message.Chat.ID, finalizedMarkdown(commands.Help.Description, botName)))

	/*
		Random command
		-------------
		Next Action: Favourite action, Unfavourite action
	*/
	case commands.Random.Key:
		// Get a random recipe
		if recipes := *recipecollection.GetRandom(1); funk.Any(recipes) {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, funk.Head(recipes).(recipecollection.Recipe).ToMessage(commands.Random.Description))

			// Add action
			message.ReplyMarkup = callback.AddActions([]int{callback.FavouriteAction, callback.UnfavouriteAction}, funk.Head(recipes).(recipecollection.Recipe).ID.Hex())
			messages = append(messages, message)
		} else {
			err = fmt.Errorf("No recipes returned for the random one")
		}

	/*
		Start command
		-------------
		Next Action: Register action
	*/
	case commands.Start.Key:
		message := tgbotapi.NewMessage(update.Message.Chat.ID, finalizedMarkdown(commands.Start.Description, botName))

		// Give option to register to new users
		if user == nil {
			message.ReplyMarkup = callback.AddActions([]int{callback.RegisterAction})
		}

		messages = append(messages, message)

	/*
		Top3 command
		-------------
		Next Action: Favourite action, Unfavourite action
	*/
	case commands.Top3.Key:
		// TODO Olli get most popular recipes

	default:
		err = fmt.Errorf("Unregocnized command %s from user [%s]", update.Message.Command(), update.Message.From.UserName)
	}

	messages = funk.Map(messages, func(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
		message.ParseMode = "Markdown"
		return message
	}).([]tgbotapi.MessageConfig)
	return &messages, err
}

// Finalize markdown with correct chars
func finalizedMarkdown(markdown string, params ...interface{}) string {
	return strings.ReplaceAll(fmt.Sprintf(markdown, params...), "''", "`")
}
