package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	raw "github.com/telegram-go-bot/go_bot/app/domain/raw_structures"
	"github.com/telegram-go-bot/go_bot/app/output/views/telegram"
)

// MessageReader - read data using telegram api
type MessageReader struct {
	bot                *tgbotapi.BotAPI
	startedLoop        bool
	receivedActivities chan *raw.Activity
}

// NewMessageReader - constructor
func NewMessageReader(botToken string) *MessageReader {
	t := MessageReader{}
	t.startedLoop = false
	t.receivedActivities = make(chan *raw.Activity, 100)

	t.bot = telegram.InitBotAPI(botToken)

	return &t
}

func (r *MessageReader) activitiesProducer() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := r.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Chat == nil {
			continue // what about new member event ? other events?
		}

		r.receivedActivities <- updateToActivity(&update)
	}

	return nil
}

// GetActivity - getActivity
func (r *MessageReader) GetActivity() (raw.Activity, error) {

	if !r.startedLoop {
		r.startedLoop = true
		go r.activitiesProducer()
	}

	activity := <-r.receivedActivities
	return *activity, nil
}

// mapping
func updateToActivity(update *tgbotapi.Update) *raw.Activity {
	if update == nil {
		return nil
	}

	var res raw.Activity

	res.Text = update.Message.Text
	res.ChatID = update.Message.Chat.ID

	return &res
}
