package activityhandlers

import (
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/input/activities"
)

// ActivityHandlerFacade - processing activities
type ActivityHandlerFacade struct {
	handlers []ICommandHandler
}

// NewActivityHandlerFacade - constructor
func NewActivityHandlerFacade(commandHandlers []ICommandHandler) *ActivityHandlerFacade {

	res := new(ActivityHandlerFacade)
	res.handlers = commandHandlers

	return res
}

// ProcessActivities - process all activities
func (h *ActivityHandlerFacade) ProcessActivities(reader activities.IActivityReader) error {
	for {
		activity, err := reader.GetActivity()
		if err != nil {
			return err
		}
		err = h.onActivity(activity)
		if err != nil {
			return err
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
