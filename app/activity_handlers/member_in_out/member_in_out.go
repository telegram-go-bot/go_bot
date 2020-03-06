package memberinout

import (
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

type impl struct {
	presenter output.IPresenter
}

// New - constructor
func New(presenter output.IPresenter) impl {
	return impl{presenter: presenter}
}

// OnHelp - display help
func (p impl) OnHelp() string {
	// NO help! Hidden cmd ;D
	return ""
}

func (p impl) OnCommand(item raw.Activity) (bool, error) {

	if item.LeftChatMember != nil {
		err := p.OnMemberLeft(item)
		return true, err
	}

	if item.NewChatMembers != nil {
		err := p.OnMemberJoined(item)
		return true, err
	}

	return false, nil
}

func (p impl) OnMemberLeft(item raw.Activity) error {
	animationID := cmn.GetOneMsgFromMany("CgADBAADwAADF8W8UOvlL6pVzc4QAg", "CgADBAADfQADFpv0UppbGdGY5UzEAg",
		"CgADAgADmQMAAlhFMUq-XYQtAqHy4QI")
	caption := cmn.GetOneMsgFromMany(
		"было приятно познакомиться", "", "так даже лучше", "F", "😢", "гуляй", "возвращайся 😢",
		"дiiiiiiiiiiiдько", "слабак", "канай отсюда", "и пусть канает")
	_, err := p.presenter.ShowGif(output.ShowAnimationData{
		ShowMessageData: output.ShowMessageData{ChatID: item.ChatID, ReplyToMsgID: item.MesssageID},
		AnimationID:     animationID, Caption: caption})
	return err
}

func (p impl) OnMemberJoined(item raw.Activity) error {
	//  TODO: move strings outside please
	animationID := cmn.GetOneMsgFromMany("CgADAgADTgADsuSgS8s5i6Vea-H9Ag")
	caption := cmn.GetOneMsgFromMany(
		"😘", "welcome",
		"Этот тут долго не продержится",
		"Поздоровайся с господами", "Shalom", "i ❤️ you",
		"И чего ты сюда пришёл? Кто тебя звал?", "Я ждал тебя",
		"Хватит добавлять сюда непонятно кого", "Рад что тебя добавили",
		"Без него было намного уютнее", "С тобой здесь будет намного уютнее",
		"Надоели люди, надоели боты", "привет, ты кто?", "привет, я бот, а ты кто?",
		"Смотрите, головка с ручками", "Лучше сразу уходи, пожалуйста, уходи",
		"😚", "🤨", "🍆🍑", "🍑", "🍆🍆")
	_, err := p.presenter.ShowGif(output.ShowAnimationData{
		ShowMessageData: output.ShowMessageData{ChatID: item.ChatID, ReplyToMsgID: item.MesssageID},
		AnimationID:     animationID, Caption: caption})

	return err
}
