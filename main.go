package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	activityhandlers "github.com/telegram-go-bot/go_bot/app/activity_handlers"
	qwertyFix "github.com/telegram-go-bot/go_bot/app/activity_handlers/fix_layout"
	"github.com/telegram-go-bot/go_bot/app/activity_handlers/goroskop"
	loopapoopa "github.com/telegram-go-bot/go_bot/app/activity_handlers/loopa_poopa"
	magicquery "github.com/telegram-go-bot/go_bot/app/activity_handlers/magic_query"
	memberinout "github.com/telegram-go-bot/go_bot/app/activity_handlers/member_in_out"
	"github.com/telegram-go-bot/go_bot/app/activity_handlers/otpetushi"
	pickfirstorsecond "github.com/telegram-go-bot/go_bot/app/activity_handlers/pick_first_or_second"
	googlephoto "github.com/telegram-go-bot/go_bot/app/activity_handlers/search_photo"
	"github.com/telegram-go-bot/go_bot/app/activity_handlers/zagadka"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	"github.com/telegram-go-bot/go_bot/app/common/vk"
	collywrapper "github.com/telegram-go-bot/go_bot/app/common/web_scrapper/colly_wrapper"
	"github.com/telegram-go-bot/go_bot/app/common/web_search/google"
	settings "github.com/telegram-go-bot/go_bot/app/database"
	in "github.com/telegram-go-bot/go_bot/app/input/activities/telegram"
	presenters "github.com/telegram-go-bot/go_bot/app/output/presenters"
	"github.com/telegram-go-bot/go_bot/app/output/views/telegram"
)

var (
	botToken             = os.Getenv("HEROKU_BOT_ID")
	baseURL              = os.Getenv("HEROKU_BASE_URL")
	dbURL                = os.Getenv("DATABASE_URL")
	botUuids             = strings.Split(os.Getenv("BOT_UIDS"), ",")
	vkLogin              = os.Getenv("VK_LOGIN")
	vkPwd                = os.Getenv("VK_PASSWORD")
	botAdminID           = os.Getenv("OWNER_ID")
	googleSearchAPIKey   = os.Getenv("GOOGLE_SEARCH_API_KEY")   // apiKey
	googleSearchEngineID = os.Getenv("GOOGLE_SEARCH_ENGINE_ID") // .cx
)

var (
	tgView       = telegram.NewTelegramAPIView(botToken)
	tgPresenter  = presenters.NewActivityPresenter(tgView)
	tgReader     = in.NewMessageReader(botToken)
	scrapper     = collywrapper.Scrapper{}
	webSearcher  = google.Init(googleSearchAPIKey, googleSearchEngineID)
	settingsInit = initSettings(dbURL)

	commandHandlers = []activityhandlers.ICommandHandler{
		pickfirstorsecond.New(tgPresenter),
		zagadka.New(tgPresenter, scrapper),
		goroskop.New(tgPresenter, scrapper),
		googlephoto.New(tgPresenter, webSearcher),
		loopapoopa.New(tgPresenter),
		memberinout.New(tgPresenter),
		otpetushi.New(tgPresenter),
		qwertyFix.New(tgPresenter),
		magicquery.New(tgPresenter, webSearcher)}
)

func initSettings(url string) error {
	cache, err := settings.NewCache(url)
	if err != nil {
		return err
	}
	settings.New(cache)
	return nil
}

func main() {
	cmn.Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

	if settingsInit != nil {
		panic("init settings failed")
	}

	vk.Init(vkLogin, vkPwd)
	tgBot := activityhandlers.New(commandHandlers, tgPresenter)
	tgBot.ProcessActivities(tgReader)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<3")
}
