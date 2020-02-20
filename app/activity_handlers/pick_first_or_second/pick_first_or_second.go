package pickfirstorsecond

import (
	"fmt"
	"regexp"

	activityhandlers "github.com/telegram-go-bot/go_bot/app/activity_handlers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

// PickFirstOrSecond -
type pickFirstOrSecond struct {
	presenter output.IPresenter
}

// New - constructor
func New(presenter output.IPresenter) pickFirstOrSecond {
	return pickFirstOrSecond{presenter: presenter}
}

// OnHelp - display help
func (p pickFirstOrSecond) OnHelp() string {
	// todo: use BOT_UIDS instead
	return "<b>!билли|billy</b> <i>smth</i> <b>или</b> <i>smth else</i> - выбираю лучший из вариантов"
}

var re = regexp.MustCompile(`!(?i)(билли|billy)\W(.+?)или ([^\?$]+)`)

func getYesNoCantPeekMsg() string {
	var items = [...]string{
		"все плохо",
		"nothing is black and white",
		"сорта говна",
		"эээ... сложно",
		"ты пидор",
		"ты билли",
		"нет"}

	itemNo := cmn.Rnd.Intn(len(items))
	return items[itemNo]
}

// OnCommand -
func (p pickFirstOrSecond) OnCommand(item activityhandlers.ActivityItem) (bool, error) {

	SendMsg := func(message string) {
		p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	}

	matches := re.FindAllStringSubmatch(item.Text, -1)
	if len(matches) == 0 {
		return false, nil
	}

	res := matches[0]

	if len(res) != 4 {
		return false,
			fmt.Errorf("PickFirstOrSecond(%s) failed to parse query. Expecting 4 values, received %d. Result of match is: %q", item.Text, len(res), res)
	}

	res = append(res[:0], res[2:]...)

	if cmn.Rnd.Intn(14) == 1 {
		// 1 times from 15 do not pick anything
		SendMsg(getYesNoCantPeekMsg())
		return true, nil
	}

	if cmn.Rnd.Intn(2) == 1 {
		SendMsg(res[0])
	} else {
		SendMsg(res[1])
	}

	return true, nil
}
