package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telegram-go-bot/go_bot/app/output"
)

// APIView - do stuff using telegram api
type APIView struct {
}

var (
	bot *tgbotapi.BotAPI
)

// InitBotAPI - InitBotAPI
func InitBotAPI(botToken string) *tgbotapi.BotAPI {
	if bot == nil {
		var err error
		bot, err = tgbotapi.NewBotAPI(botToken)
		if err != nil {
			log.Panic(err)
		}
	}
	return bot
}

// NewTelegramAPIView - constructor
func NewTelegramAPIView(botToken string) *APIView {
	t := new(APIView)
	InitBotAPI(botToken)

	return t
}

// ShowMessage - display msg using telegram bot api
func (t *APIView) ShowMessage(msg output.ViewMessageData) error {

	msgOut := tgbotapi.NewMessage(msg.ChatID, msg.Text)
	bot.Send(msgOut)

	return nil
}
