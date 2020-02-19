package collywrapper

import (
	"github.com/gocolly/colly"
	iscrapper "github.com/telegram-go-bot/go_bot/app/common/web_scrapper"
)

// Scrapper - gocolly wrapper
type Scrapper struct {
	collector *colly.Collector
}

// Init - initialize
func (c Scrapper) Init() iscrapper.Interface {
	var s Scrapper
	s.collector = colly.NewCollector()
	if s.collector == nil {
		panic("error initializing gocolly")
	}
	return s
}

// OnHTML -
func (c Scrapper) OnHTML(selector string, cb iscrapper.OnHTMLCallback) error {

	// On every a element which has href attribute call callback
	c.collector.OnHTML(selector, func(e *colly.HTMLElement) {
		cb(e.Text)
	})

	return nil
}

// Visit -
func (c Scrapper) Visit(url string) error {
	return c.collector.Visit(url)
}
