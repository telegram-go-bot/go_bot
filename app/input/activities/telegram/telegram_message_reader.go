package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	settings "github.com/telegram-go-bot/go_bot/app/database"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/output/views/telegram"
)

// MessageReader - read data using telegram api
type MessageReader struct {
	bot                *tgbotapi.BotAPI
	startedLoop        bool
	receivedActivities chan *raw.Activity
}

// returns false if no message should be processed from definite chat
// Message != nil
// Message.Chat != nil
func confirureSettingsForMessage(msg *tgbotapi.Message) bool {
	chatID := msg.Chat.ID
	chatInfo, err := settings.Inst().GetChatInfo(chatID)
	if err != nil {
		// ok no chatInfo exists - create new one
		chatInfo = &raw.ChatInfo{}
		chatInfo.ChatID = chatID
		chatInfo.Name = msg.Chat.Title
		settings.Inst().AddRecord(chatInfo)
		return true
	}
	return chatInfo.Enabled
}

func configureUsersForMessage(msg *tgbotapi.Message) {
	if msg.From == nil {
		return
	}
	_, err := settings.Inst().GetChatUser(msg.From.ID)
	if err != nil { // no user exists. persist one
		newUser := messageToChatUser(msg)
		settings.Inst().AddRecord(newUser)
	}
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

		// skip chat if it is disabled, or add if new
		if !confirureSettingsForMessage(update.Message) {
			continue
		}

		// process users
		configureUsersForMessage(update.Message)

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

func toRawUserPtr(user *tgbotapi.User) *raw.User {
	if user == nil {
		return nil
	}

	res := toRawUser(user)
	return &res
}

func toRawUser(user *tgbotapi.User) raw.User {
	return raw.User{
		ID:       user.ID,
		UserName: user.UserName,
		IsBot:    user.IsBot}
}

func messageToRawActivity(msg *tgbotapi.Message, activityOut *raw.Activity) {
	if msg == nil {
		return
	}

	if msg.Chat == nil {
		log.Printf("Received message with empty ChatID !!! Message text is: %s\n", msg.Text)
		return
	}

	activityOut.Text = msg.Text
	activityOut.ChatID = msg.Chat.ID
	activityOut.Command = msg.Command()
	activityOut.MesssageID = msg.MessageID
}

func updateLeftChatMember(activity *raw.Activity, message *tgbotapi.Message) {
	if message.LeftChatMember != nil {
		leftUser := toRawUser(message.LeftChatMember)
		activity.LeftChatMember = &leftUser
	}
}

func updateNewChatMembers(activity *raw.Activity, message *tgbotapi.Message) {
	if message.NewChatMembers == nil {
		return
	}
	countOfNewMembers := len(*message.NewChatMembers)
	if countOfNewMembers == 0 {
		return
	}

	newMembers := make([]raw.User, countOfNewMembers)
	activity.NewChatMembers = &newMembers

	for i := 0; i < countOfNewMembers; i++ {
		(*activity.NewChatMembers)[i] = toRawUser(&(*message.NewChatMembers)[i])
	}
}

// mapping
func updateToActivity(update *tgbotapi.Update) *raw.Activity {
	if update == nil {
		return nil
	}

	var activity raw.Activity
	if update.Message != nil {
		messageToRawActivity(update.Message, &activity)
		updateLeftChatMember(&activity, update.Message)
		updateNewChatMembers(&activity, update.Message)
		if update.Message.ReplyToMessage != nil {
			activity.RepliedTo = new(raw.Activity)
			messageToRawActivity(update.Message.ReplyToMessage, activity.RepliedTo)
		}
	}

	return &activity
}

func messageToChatUser(msg *tgbotapi.Message) *raw.ChatUser {
	usr := msg.From
	newUser := &raw.ChatUser{
		UserID:       usr.ID,
		UserName:     usr.UserName,
		ChatID:       msg.Chat.ID,
		ChatUserInfo: raw.ChatUserInfo{FirstName: usr.FirstName, LastName: usr.LastName}}
	return newUser
}
