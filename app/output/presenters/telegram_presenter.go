package activities

import (
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
	var msgData output.ViewMessageData
	msgData.Text = data.Text
	msgData.ChatID = data.ChatID
	msgData.ReplyToMsgID = data.ReplyToMsgID

	msgID, err := s.view.ShowMessage(msgData)
	if err != nil {
		return 0, err
	}

	return msgID, nil
}
