package goroskop

import (
	"strconv"
	"sync"
	"time"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	webscrapper "github.com/telegram-go-bot/go_bot/app/common/web_scrapper"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/output"
)

type cachedItem struct {
	dateVal int
	msg     string
}

type impl struct {
	presenter         output.IPresenter
	scrapper          webscrapper.Interface
	horoscopeForToday map[int64]cachedItem // key - ChatID, val - cached horoscope message
	cacheLock         *sync.RWMutex
}

// New - constructor
func New(presenter output.IPresenter, scrapper webscrapper.Interface) impl {
	return impl{
		presenter:         presenter,
		scrapper:          scrapper,
		horoscopeForToday: make(map[int64]cachedItem),
		cacheLock:         &sync.RWMutex{}}
}

// OnHelp - display help
func (p impl) OnHelp() string {
	return "<b>!чобудет|шобудет</b> - расскажу что тебя ждет"
}

func (p impl) reuseHoroscopeIfExists(chatID int64) (string, bool) {
	p.cacheLock.RLock()
	todays, exists := p.horoscopeForToday[chatID]
	p.cacheLock.RUnlock()
	if !exists {
		return "", false
	}

	_, _, todaysDay := time.Now().Date()
	if todaysDay != todays.dateVal {
		return "", false
	}

	prefix := cmn.GetOneMsgFromMany("уже было: ", "та говорил же уже: ", "еще раз: ", "", "", "")

	return prefix + todays.msg, true
}

func (p impl) getNewHoroscopeForToday(chatID int64) (string, error) {

	scrapper := p.scrapper.Init()
	visited := make([]string, 0, 5)

	err := scrapper.OnHTML("P", func(text string) {
		if len(text) != 0 {
			visited = append(visited, text)
		}
	})
	if err != nil {
		return "", err
	}

	err = scrapper.Visit("http://stoboi.ru/gorodaily/horoscope.php?id=" + strconv.Itoa(cmn.Rnd.Intn(12)+1))
	if err != nil {
		return "", err
	}

	itemNo := cmn.Rnd.Intn(len(visited))

	msg := visited[itemNo]

	_, _, todaysDay := time.Now().Date()
	p.cacheLock.Lock()
	p.horoscopeForToday[chatID] = cachedItem{dateVal: todaysDay, msg: msg}
	p.cacheLock.Unlock()

	return msg, nil
}

// OnCommand -
func (p impl) OnCommand(item raw.Activity) (bool, error) {

	_, isThisCommand := helpers.IsOnCommand(item.Text, []string{"чобудет", "шобудет"})
	if !isThisCommand {
		return false, nil
	}

	SendMsg := func(message string) (int, error) {
		return p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	}

	// if we've cached horoscope for today already - reuse it
	msg, exists := p.reuseHoroscopeIfExists(item.ChatID)
	if exists {
		SendMsg(msg)
		return true, nil
	}

	msg, err := p.getNewHoroscopeForToday(item.ChatID)
	if err != nil {
		return false, err
	}

	SendMsg(msg)
	return true, nil
}
