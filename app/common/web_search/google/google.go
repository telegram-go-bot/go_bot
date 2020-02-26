package google

import (
	"log"
	"net/http"
	"path/filepath"

	isearcher "github.com/telegram-go-bot/go_bot/app/common/web_search"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

const (
	maxImagesResult = 10
)

// Search - init and use search engine.
type Search struct {
	apiKey string // search_api_key
	cx     string // search_engine_id
}

// Init new search
func Init(apiKey string, searchEngineID string) isearcher.Searcher {
	var s Search
	s.apiKey = apiKey
	s.cx = searchEngineID
	return s
}

// MaxResponseItemsCount - max number of images that could be returned by single query
func (s Search) MaxResponseItemsCount() int {
	return maxImagesResult
}

/*//SearchImageURL - returns random picked image URL among as much as @gMaxImagesResult results
func SearchImageURL(query string) string {
	images := searchImageMaxImagesResultDepth(query)
	if len(images) == 0 {
		return ""
	}
	pickN := rand.Intn(len(images))
	return images[pickN]
}

//SearchImage10Depth - queries for the first @MaxImagesResult images only
func searchImageMaxImagesResultDepth(query string) []string {
	return SearchImage(query, maxImagesResult)
}*/

//SearchImage - search for a query string, returns as max as searchDepth value
func (s Search) SearchImage(query string, searchDepth int) []string {

	log.Printf("Searching for: %s\n", query)

	hc := &http.Client{Transport: &transport.APIKey{Key: s.apiKey}}
	svc, err := customsearch.New(hc)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := svc.Cse.List(query).
		Cx(s.cx).
		SearchType("image").
		Do()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(s.cx)
	log.Printf("Google response contains (%d) Items...\n", len(resp.Items))

	var res []string
	for _, result := range resp.Items {

		if result.Image == nil {
			continue
		}

		// dont process image urls like x-raw-image:///239d4d29553e6d...
		if len(filepath.Ext(result.Link)) == 0 {
			continue
		}

		res = append(res, result.Link)
		log.Printf("\t%s\n", result.Link)

		if len(res) > searchDepth {
			break
		}
	}
	return res
}
