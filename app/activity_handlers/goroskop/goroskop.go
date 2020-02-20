package goroskop

import (
	"strconv"

	activityhandlers "github.com/telegram-go-bot/go_bot/app/activity_handlers"
	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	webscrapper "github.com/telegram-go-bot/go_bot/app/common/web_scrapper"
	"github.com/telegram-go-bot/go_bot/app/output"
)

type impl struct {
	presenter output.IPresenter
	scrapper  webscrapper.Interface
}

// New - constructor
func New(presenter output.IPresenter, scrapper webscrapper.Interface) impl {
	return impl{presenter: presenter, scrapper: scrapper}
}

// OnHelp - display help
func (p impl) OnHelp() string {
	return "<b>!чобудет|шобудет</b> - расскажу что тебя ждет"
}

// OnCommand -
func (p impl) OnCommand(item activityhandlers.ActivityItem) (bool, error) {

	_, isThisCommand := helpers.IsOnCommand(item.Text, []string{"чобудет", "шобудет"})
	if !isThisCommand {
		return false, nil
	}

	scrapper := p.scrapper.Init()

	visited := make([]string, 0, 5)

	err := scrapper.OnHTML("P", func(text string) {
		if len(text) != 0 {
			visited = append(visited, text)
		}
	})
	if err != nil {
		return false, err
	}

	err = scrapper.Visit("http://stoboi.ru/gorodaily/horoscope.php?id=" + strconv.Itoa(cmn.Rnd.Intn(12)+1))
	if err != nil {
		return false, err
	}

	itemNo := cmn.Rnd.Intn(len(visited))

	SendMsg := func(message string) (int, error) {
		return p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	}

	SendMsg(visited[itemNo])
	return true, nil
}
