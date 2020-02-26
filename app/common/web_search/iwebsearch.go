package websearch

// Searcher - search for image url, or some text
type Searcher interface {
	MaxResponseItemsCount() int
	//SearchImage - search for a query string, returns list of image URLs max as searchDepth value
	SearchImage(query string, searchDepth int) []string
}
