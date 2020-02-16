package activityhandlers

// ActivityItem - representation of activity, but in terms of commandhandler
type ActivityItem struct {
	Text   string
	ChatID int64
}

// ICommandHandler - common interface for all handlers
type ICommandHandler interface {
	OnHelp() string
	OnCommand(ActivityItem) (bool, error)
}
