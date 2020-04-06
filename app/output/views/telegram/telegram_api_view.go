package telegram

import (
	"errors"
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
func (t *APIView) ShowMessage(msg output.ViewMessageData) (int, error) {

	if len(msg.Text) == 0 {
		return 0,
			errors.New("<TgApiView::ShowMessage> failed to show message - it is empty")
	}

	msgOut := tgbotapi.NewMessage(msg.ChatID, msg.Text)
	if msg.ReplyToMsgID != 0 {
		msgOut.ReplyToMessageID = msg.ReplyToMsgID
	}

	if msg.ParseMode == output.ParseModeHTML {
		msgOut.ParseMode = tgbotapi.ModeHTML
	}

	sent, err := bot.Send(msgOut)
	if err != nil {
		return 0, err
	}

	return sent.MessageID, nil
}

// ShowImage - display image using telegram bot api
func (t *APIView) ShowImage(msg output.ViewImageData) (int, error) {

	if len(msg.ImageData) == 0 {
		return 0,
			errors.New("<TgApiView::ShowImage> image buffer is empty")
	}

	fileBytes := tgbotapi.FileBytes{Name: "", Bytes: msg.ImageData}
	photoMsg := tgbotapi.NewPhotoUpload(msg.ChatID, fileBytes)
	
	if msg.ReplyToMsgID != 0 {
		photoMsg.ReplyToMessageID = msg.ReplyToMsgID
	}
	
	sent, err := bot.Send(photoMsg)
	if err != nil {
		return 0, err
	}

	return sent.MessageID, nil
}

// ShowAnimation - display image using telegram bot api
func (t *APIView) ShowAnimation(animation output.ViewAnimationData) (int, error) {
	if len(animation.AnimationID) == 0 {
		return 0,
			errors.New("<TgApiView::ShowAnimation> empty image animation iD")
	}

	msg := tgbotapi.NewAnimationShare(animation.ChatID, animation.AnimationID)
	if animation.ReplyToMsgID != 0 {
		msg.ReplyToMessageID = animation.ReplyToMsgID
	}

	if len(animation.Caption) > 0 {
		msg.Caption = animation.Caption
	}

	sent, err := bot.Send(msg)
	if err != nil {
		return 0, err
	}

	return sent.MessageID, nil
}
