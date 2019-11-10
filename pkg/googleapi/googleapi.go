package googleapi

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Recipe fron google API
type Recipe struct {
	Rating float32
	Title  string
	URL    string
}

// SearchRecipes from google API
func SearchRecipes(recipe string) (recipes *[]Recipe) {

	// Get HTML from google here
	html, err := getGoogleSearchResult(recipe)
	if err != nil {
		return &[]Recipe{}
	}
	// Parse results to SearchResults
	return parseGoogleRecipes(html)
}

// Get google search result
func getGoogleSearchResult(query string) (html *string, err error) {

	response, err := http.Get("https://www.devdungeon.com/")
	if err != nil {
		log.Printf("Failed to get search results from google. %s", err)
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to parse search results from google. %s", err)
	}

	stringBody := string(body)
	return &stringBody, err
}

// Parse google results html page
func parseGoogleRecipes(html *string) (recipes *[]Recipe) {
	return &[]Recipe{}
}
