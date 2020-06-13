package covid

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

// Impl - impl
type Impl struct {
	presenter output.IPresenter
}

// New - constructor
func New(presenter output.IPresenter) Impl {
	return Impl{presenter: presenter}
}

// OnHelp - display help
func (p Impl) OnHelp() string {
	return "<b>!covid|ковид</b> - свежая статистика по ковиду"
}

type covidCountry struct {
	Country        string
	CountryCode    string
	Date           string
	NewConfirmed   int64
	NewDeaths      int64
	NewRecovered   int64
	Slug           string
	TotalConfirmed int64
	TotalDeaths    int64
	TotalRecovered int64
}

type covidJSON struct {
	Countries []covidCountry
	Date      string
	Global    struct {
		NewConfirmed   int64
		NewDeaths      int64
		NewRecovered   int64
		TotalConfirmed int64
		TotalDeaths    int64
		TotalRecovered int64
	}
}

func getCountryByCode(allCountries []covidCountry, codeToSearchFor string) int {
	for id, country := range allCountries {
		if country.CountryCode == codeToSearchFor {
			return id
		}
	}
	return -1
}

func countryFromASCII(ascii string) string {
	r := []rune(ascii)
	for i := 0; i < len(r); i++ {
		r[i] += 127397
	}
	return string(r)
}

func toString(num int64) string {
	s := strconv.FormatInt(num, 10)
	const devider int = 3
	if len(s) <= devider {
		return s
	}

	var buffer bytes.Buffer
	counter := len(s) % devider

	for id, ch := range s {
		if id < counter {
			buffer.WriteRune(ch)
			continue
		}
		if id != 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteRune(ch)
		counter += devider
	}

	return buffer.String()
}

func formatResponse(item raw.Activity, all covidJSON, ua covidCountry) (output.ShowMessageData, error) {
	res := output.ShowMessageData{ChatID: item.ChatID}

	var buffer bytes.Buffer
	buffer.WriteString("Заразилось: *" + toString(all.Global.TotalConfirmed) + "* (+" + toString(all.Global.NewConfirmed) + ")\n")
	buffer.WriteString("Излечилось: *" + toString(all.Global.TotalRecovered) + "* (+" + toString(all.Global.NewRecovered) + ")\n")
	buffer.WriteString("Умерло: *" + toString(all.Global.TotalDeaths) + "* (+" + toString(all.Global.NewDeaths) + ")\n")
	buffer.WriteString("\n")

	for _, country := range all.Countries {
		buffer.WriteString(countryFromASCII(country.CountryCode) + " *" + toString(country.TotalConfirmed) + "* (+" + toString(country.NewConfirmed) + ")\n")
	}

	buffer.WriteString("\n")
	buffer.WriteString(countryFromASCII(ua.CountryCode) + "\n")
	buffer.WriteString("Заразилось: *" + toString(ua.TotalConfirmed) + "* (+" + toString(ua.NewConfirmed) + ")\n")
	buffer.WriteString("Излечилось: *" + toString(ua.TotalRecovered) + "* (+" + toString(ua.NewRecovered) + ")\n")
	buffer.WriteString("Умерло: *" + toString(ua.TotalDeaths) + "* (+" + toString(ua.NewDeaths) + ")\n")

	res.Text = buffer.String()
	res.ParseMode = output.ParseModeMarkdown

	return res, nil
}

// OnCommand -
func (p Impl) OnCommand(item raw.Activity) (bool, error) {

	_, isThisCommand := helpers.IsOnCommand(item.Text, []string{"ковид", "covid", "корона", "corona"})
	if !isThisCommand {
		return false, nil
	}

	url := "https://api.covid19api.com/summary"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	//req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	summary := covidJSON{}
	jsonErr := json.Unmarshal(body, &summary)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	sort.Slice(summary.Countries, func(i, j int) bool {
		return summary.Countries[i].TotalConfirmed > summary.Countries[j].TotalConfirmed
	})

	// get Ukraine
	uaID := getCountryByCode(summary.Countries, "UA")
	if uaID == -1 {
		return true, errors.New("cant find UA :O")
	}
	ua := summary.Countries[uaID]

	summary.Countries = summary.Countries[:5]

	resp, err := formatResponse(item, summary, ua)
	if err != nil {
		return true, errors.New("eek cant format response")
	}

	//p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	p.presenter.ShowMessage(resp)

	return true, nil
}
