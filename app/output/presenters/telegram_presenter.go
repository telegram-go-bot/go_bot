package activities

import (
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	"github.com/telegram-go-bot/go_bot/app/output"
)

// ActivityPresenter - produces ViewMessageData from msg text
type ActivityPresenter struct {
	view output.IView
}

// NewActivityPresenter - constructor for ActivityPresenter
func NewActivityPresenter(outputView output.IView) *ActivityPresenter {
	res := new(ActivityPresenter)
	res.view = outputView
	return res
}

// ShowMessage - display dummy message
func (s ActivityPresenter) ShowMessage(data output.ShowMessageData) (int, error) {
	var msg output.ViewMessageData
	msg.Text = data.Text
	msg.ChatID = data.ChatID
	msg.ParseMode = data.ParseMode
	msg.ReplyToMsgID = data.ReplyToMsgID

	msgID, err := s.view.ShowMessage(msg)
	if err != nil {
		return 0, err
	}

	return msgID, nil
}

// ShowImage - Download image from URL and show it
func (s ActivityPresenter) ShowImage(imageData output.ShowImageData) (int, error) {

	var msg output.ViewImageData

	if imageData.RawImageData != nil {
		msg.ImageData = imageData.RawImageData
	} else {
		var err error
		msg.ImageData, err = cmn.DownloadFileByURL(imageData.ImageURL)
		if err != nil {
			return 0, err
		}
	}

	msg.Text = imageData.Text
	msg.ChatID = imageData.ChatID
	msg.ReplyToMsgID = imageData.ReplyToMsgID
	msg.ParseMode = imageData.ParseMode

	msgID, err := s.view.ShowImage(msg)
	if err != nil {
		return 0, err
	}

	return msgID, nil
}

// ShowGif - display animation
func (s ActivityPresenter) ShowGif(animationData output.ShowAnimationData) (int, error) {
	var data output.ViewAnimationData
	data.Text = animationData.Text
	data.ChatID = animationData.ChatID
	data.ReplyToMsgID = animationData.ReplyToMsgID
	data.ParseMode = animationData.ParseMode
	data.AnimationID = animationData.AnimationID
	data.Caption = animationData.Caption

	msgID, err := s.view.ShowAnimation(data)
	if err != nil {
		return 0, err
	}

	return msgID, nil
}
