package cmn

import raw "github.com/telegram-go-bot/go_bot/app/domain"

// GetReplyToMessageID - activity - a message that triggers bot. E.g. !command
func GetReplyToMessageID(activity raw.Activity) int {
	if activity.RepliedTo != nil {
		return activity.RepliedTo.MesssageID
	}
	return activity.MesssageID
}
