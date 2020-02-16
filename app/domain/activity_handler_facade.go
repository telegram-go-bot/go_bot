package domain

import (
	activityhandlers "github.com/telegram-go-bot/go_bot/app/activity_handlers"
	raw "github.com/telegram-go-bot/go_bot/app/domain/raw_structures"
	"github.com/telegram-go-bot/go_bot/app/input/activities"
)

// ActivityHandlerFacade - processing activities
type ActivityHandlerFacade struct {
	handlers []activityhandlers.ICommandHandler
}

// NewActivityHandlerFacade - constructor
func NewActivityHandlerFacade(commandHandlers []activityhandlers.ICommandHandler) *ActivityHandlerFacade {

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
		err = h.onActivity(activityToActivityItem(activity))
		if err != nil {
			return err
		}
	}
}

func (h *ActivityHandlerFacade) onActivity(activity activityhandlers.ActivityItem) error {

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

// mapping
func activityToActivityItem(activity raw.Activity) activityhandlers.ActivityItem {
	var res activityhandlers.ActivityItem
	res.Text = activity.Text
	res.ChatID = activity.ChatID
	return res
}
