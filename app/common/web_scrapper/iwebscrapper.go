package webscrapper

// OnHTMLCallback -
type OnHTMLCallback func(text string)

// Interface - get data from web pages
type Interface interface {
	Init() Interface
	Visit(url string) error
	OnHTML(selector string, cb OnHTMLCallback) error
}
