package zagadka

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	webscrapper "github.com/telegram-go-bot/go_bot/app/common/web_scrapper"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

// todo: https://www.anekdotovmir.ru/shutki/zagadki-shutochnye/
// 		 https://www.anekdotovmir.ru/shutki/zagadki-shutochnye/shutki-zadachi/
//       https://crossword.nalench.com/others/37330-poshlye-zagadki.html
//		 https://crossword.nalench.com/others/37525-smeshnye-zagadki-dlja-vzroslyh.html

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

type lastUsed struct {
	page    int
	itemNum int
}

var questsLastUsedInfo = map[int64] /*chatID*/ lastUsed{}

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

	questLastUsedInfo, found := questsLastUsedInfo[item.ChatID]
	if !found {
		questLastUsedInfo = lastUsed{page: 1, itemNum: 19}
		questsLastUsedInfo[item.ChatID] = questLastUsedInfo
	}

	pageNo := questLastUsedInfo.page
	itemNo := questLastUsedInfo.itemNum

	log.Printf("Zagadka. page: %d, item: %d", pageNo, itemNo)

	itemNo++
	if pageNo*20 < itemNo || itemNo > 130 {
		// inc page
		if pageNo == 7 {
			pageNo = 1
			itemNo = 1
		} else {
			pageNo++
			itemNo = (pageNo-1)*20 + 1
		}
	}

	localItemNoOnPage := itemNo - (pageNo-1)*20

	items := make([]questItem, 0, 20)

	scrapper := p.scrapper.Init()

	scrapper.OnHTML("div.ttexts", func(text string) {
		addNewQuest(text, &items)
	})

	scrapper.Visit("http://allriddles.ru/ru/riddles/joke/p" + strconv.Itoa(pageNo) + "/")

	if len(items) == 0 {
		SendMsg(cmn.GetFailMsg())
		return
	}

	itemID := localItemNoOnPage - 1

	sentID, err := SendMsg(items[itemID].zagadka)
	if err != nil {
		SendMsg(cmn.GetFailMsg())
	}

	// update reply to msg id
	quest := items[itemID]
	quest.msgID = sentID
	activeItems[item.ChatID] = quest

	questsLastUsedInfo[item.ChatID] = lastUsed{page: pageNo, itemNum: itemNo}
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
