package covid

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

// Impl - impl
type Impl struct {
	presenter output.IPresenter
	ukraine   covidCountry

	currentCovidData covidData
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
	title               string
	code                string
	totalCases          int64
	totalRecovered      int64
	totalUnresolved     int64
	totalDeaths         int64
	totalActiveCases    int64
	totalSeriousCases   int64
	totalNewCasesToday  int64
	totalNewDeathsToday int64
}

type covidData struct {
	covidCountries list.List
	Date           *time.Time
	Global         struct {
		NewConfirmed   int64
		NewDeaths      int64
		TotalConfirmed int64
		TotalDeaths    int64
	}
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
	if s[0] == '-' {
		s = s[1:]
		buffer.WriteString("-")
	}

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

func formatResponse(item raw.Activity, all covidData, ua covidCountry) (output.ShowMessageData, error) {
	res := output.ShowMessageData{ChatID: item.ChatID}

	var buffer bytes.Buffer
	buffer.WriteString("Заразилось: *" + toString(all.Global.TotalConfirmed) + "* (+" + toString(all.Global.NewConfirmed) + ")\n")
	//buffer.WriteString("Излечилось: *" + toString(all.Global.TotalRecovered) + "* (+" + toString(all.Global.NewRecovered) + ")\n")
	buffer.WriteString("Умерло: *" + toString(all.Global.TotalDeaths) + "* (+" + toString(all.Global.NewDeaths) + ")\n")
	buffer.WriteString("\n")

	for e := all.covidCountries.Front(); e != nil; e = e.Next() {
		country := e.Value.(covidCountry)
		buffer.WriteString(countryFromASCII(country.code) + " *" + toString(country.totalCases) + "* (+" + toString(country.totalNewCasesToday) + ")\n")
	}

	buffer.WriteString("\n")
	buffer.WriteString(countryFromASCII(ua.code) + "\n")
	buffer.WriteString("Заразилось: *" + toString(ua.totalCases) + "* (+" + toString(ua.totalNewCasesToday) + ")\n")
	buffer.WriteString("Активных: *" + toString(ua.totalActiveCases+ua.totalSeriousCases) + "* (+" + toString(ua.totalNewCasesToday) + ")\n")
	buffer.WriteString("Умерло: *" + toString(ua.totalDeaths) + "* (+" + toString(ua.totalNewDeathsToday) + ")\n")

	res.Text = buffer.String()
	res.ParseMode = output.ParseModeMarkdown

	return res, nil
}

func (p Impl) addCountry(countries *list.List, currentCovidData *covidData, new covidCountry, maxCountries int) {

	currentCovidData.Global.NewConfirmed += new.totalNewCasesToday
	currentCovidData.Global.NewDeaths += new.totalNewDeathsToday
	currentCovidData.Global.TotalConfirmed += new.totalCases
	currentCovidData.Global.TotalDeaths += new.totalDeaths

	if countries.Len() == 0 {
		countries.PushBack(new)
		return
	}

	for e := countries.Front(); e != nil; e = e.Next() {
		if new.totalNewCasesToday >= e.Value.(covidCountry).totalNewCasesToday {
			countries.InsertBefore(new, e)
			break
		}
	}

	if countries.Len() > maxCountries {
		countries.Remove(countries.Back())
	}
}

func (p Impl) interfaceToCountry(countryInterface interface{}) (covidCountry, *covidCountry, error) {
	var res covidCountry
	for itemKey, itemValue := range countryInterface.(map[string]interface{}) {
		switch itemKey {
		case "title":
			res.title = itemValue.(string)
		case "code":
			res.code = itemValue.(string)
		case "total_cases":
			res.totalCases = int64(itemValue.(float64))
		case "total_recovered":
			res.totalRecovered = int64(itemValue.(float64))
		case "total_unresolved":
			res.totalUnresolved = int64(itemValue.(float64))
		case "total_deaths":
			res.totalDeaths = int64(itemValue.(float64))
		case "total_active_cases":
			res.totalActiveCases = int64(itemValue.(float64))
		case "total_serious_cases":
			res.totalSeriousCases = int64(itemValue.(float64))
		case "total_new_cases_today":
			res.totalNewCasesToday = int64(itemValue.(float64))
		case "total_new_deaths_today":
			res.totalNewDeathsToday = int64(itemValue.(float64))
		}
	}

	var ukraine *covidCountry

	if res.code == "UA" {
		ukraine = &res
	}

	return res, ukraine, nil
}

// OnCommand -
func (p Impl) OnCommand(item raw.Activity) (bool, error) {

	_, isThisCommand := helpers.IsOnCommand(item.Text, []string{"ковид", "covid", "корона", "corona", "насморк"})
	if !isThisCommand {
		return false, nil
	}

	p.currentCovidData = covidData{}

	url := "https://api.thevirustracker.com/free-api?countryTotals=ALL"

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Maximum of 2 secs
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

	var summary interface{}
	jsonErr := json.Unmarshal(body, &summary)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	itemsMap := summary.(map[string]interface{})

	// Loop through the Items; we're not interested in the key, just the values
	for key, v := range itemsMap {

		if key != "countryitems" {
			continue
		}
		// Use type assertions to ensure that the value's a JSON object
		switch v.(type) {
		// The value is an Item, represented as a generic interface
		case interface{}:

			datas := v.([]interface{})

			for _, data := range datas {
				for _, itemValue := range data.(map[string]interface{}) {
					switch contriesIbj := itemValue.(type) {
					case string:
						log.Printf("%+v", contriesIbj)
					case interface{}:
						country, ukr, err := p.interfaceToCountry(contriesIbj)
						if err == nil {
							if ukr != nil {
								p.ukraine = *ukr
							}
							p.addCountry(&p.currentCovidData.covidCountries, &p.currentCovidData, country, 5)
						}
					}

				}
			}
		// Not a JSON object; handle the error
		default:
			fmt.Println("Expecting a JSON object; got something else")
		}
	}

	resp, err := formatResponse(item, p.currentCovidData, p.ukraine)
	if err != nil {
		return true, errors.New("eek cant format response")
	}

	//p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	p.presenter.ShowMessage(resp)

	return true, nil
}
