package zagadka

import (
	"regexp"
	"strconv"
	"strings"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	webscrapper "github.com/telegram-go-bot/go_bot/app/common/web_scrapper"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

type zagadka struct {
	presenter output.IPresenter
	scrapper  webscrapper.Interface
}

// New - constructor
func New(presenter output.IPresenter, scrapper webscrapper.Interface) zagadka {
	return zagadka{presenter: presenter, scrapper: scrapper}
}

// OnHelp - display help
func (p zagadka) OnHelp() string {
	return "<b>!загадка</b> затем <b>!отгадка|ответ|разгадка</b> - загадки для выпускников Гарварда"
}

// OnCommand -
func (p zagadka) OnCommand(item raw.Activity) (bool, error) {
	var isZagadka, isRazgadka bool

	_, isZagadka = helpers.IsOnCommand(item.Text, []string{"загадка"})
	if !isZagadka {
		_, isRazgadka = helpers.IsOnCommand(item.Text, []string{"отгадка", "ответ", "разгадка"})
		if !isRazgadka {
			return false, nil
		}
	}

	if isZagadka {
		go p.onZagadka(item)
	} else {
		go p.onOtgadka(item)
	}

	return true, nil
}

func (p zagadka) onZagadka(item raw.Activity) {

	SendMsg := func(message string) (int, error) {
		return p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	}

	items := make([]questItem, 0, 20)

	scrapper := p.scrapper.Init()

	scrapper.OnHTML("div.ttexts", func(text string) {
		addNewQuest(text, &items)
	})

	scrapper.Visit("http://allriddles.ru/ru/riddles/joke/p" + strconv.Itoa(cmn.Rnd.Intn(7)+1) + "/")

	if len(items) == 0 {
		SendMsg(cmn.GetFailMsg())
		return
	}

	itemID := cmn.Rnd.Intn(len(items) + 1)

	sentID, err := SendMsg(items[itemID].zagadka)
	if err != nil {
		SendMsg(cmn.GetFailMsg())
	}

	// update reply to msg id
	quest := items[itemID]
	quest.msgID = sentID
	activeItems[item.ChatID] = quest
}

func (p zagadka) onOtgadka(item raw.Activity) {

	SendMsgTo := func(message string, replyTo int) (int, error) {
		return p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message, ReplyToMsgID: replyTo})
	}

	quest, exists := activeItems[item.ChatID]
	if !exists || len(quest.otvet) == 0 {
		return
	}

	SendMsgTo(quest.otvet, quest.msgID)

	// remove quest, no longer used
	delete(activeItems, item.ChatID)
}

//------------------------------------------------------------------------------------

type questItem struct {
	zagadka string
	otvet   string
	msgID   int
}

type questItems = map[string]string

var (
	activeItems = make(map[int64]questItem)
	questRegex  = regexp.MustCompile(`([^\[]+)\[ Ответ \]([^\[]+)`)
)

func addNewQuest(rawString string, item *[]questItem) {
	s1 := questRegex.FindStringSubmatch(rawString)
	if len(s1) != 3 {
		return
	}
	zagadka := strings.TrimSpace(s1[1])
	otvet := strings.TrimSpace(s1[2])

	*item = append(*item, questItem{zagadka: zagadka, otvet: otvet})
}
