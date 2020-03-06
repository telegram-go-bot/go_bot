package activityhandlers

import (
	"log"
	"strings"

	raw "github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/input/activities"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

// ActivityHandlerFacade - processing activities
type ActivityHandlerFacade struct {
	handlers  []ICommandHandler
	helpMsg   string
	presenter output.IPresenter
}

// New - constructor
func New(commandHandlers []ICommandHandler, presenter output.IPresenter) *ActivityHandlerFacade {

	res := new(ActivityHandlerFacade)
	res.handlers = commandHandlers
	res.helpMsg = res.initHelpMsg()
	res.presenter = presenter

	return res
}

func (h *ActivityHandlerFacade) initHelpMsg() string {
	var htmlT strings.Builder

	var endl = func() {
		htmlT.WriteByte(10)
		htmlT.WriteByte(9)
	}

	htmlT.WriteString(`Список команд:`)
	endl()

	for _, handler := range h.handlers {
		htmlT.WriteString(handler.OnHelp())
		endl()
	}

	return htmlT.String()
}

// ProcessActivities - process all activities
func (h *ActivityHandlerFacade) ProcessActivities(reader activities.IActivityReader) error {
	for {
		activity, err := reader.GetActivity()
		if err != nil {
			return err
		}

		if activity.Command == "help" {
			h.presenter.ShowMessage(
				output.ShowMessageData{
					ChatID:    activity.ChatID,
					Text:      h.helpMsg,
					ParseMode: output.ParseModeHTML})
			continue
		}

		err = h.onActivity(activity)
		if err != nil {
			log.Printf("ProcessActivity error: " + err.Error())
			continue
		}
	}
}

func (h *ActivityHandlerFacade) onActivity(activity raw.Activity) error {

	for _, handler := range h.handlers {
		done, err := handler.OnCommand(activity)
		if err != nil {
			return err
		}
		// stop when one capture succeed ?
		if done {
			break
		}
	}

	return nil
}
