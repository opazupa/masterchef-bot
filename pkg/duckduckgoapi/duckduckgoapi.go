package duckduckgoapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Recipe fron duckduckgo API
type Recipe struct {
	Rating float32
	Title  string
	URL    string
}

const (
	duckDuckGoAPI = "https://api.duckduckgo.com"
)

// SearchRecipes from duckduckgo API
func SearchRecipes(recipe string) (recipes *[]Recipe) {

	// Get search ressults from duckduckgo
	html, err := getDuckDuckGoSearchResult(recipe)
	if err != nil {
		return &[]Recipe{}
	}
	// Parse results to SearchResults
	return parseDuckDuckGoRecipes(html)
}

// Get duckduckgo search result
func getDuckDuckGoSearchResult(query string) (html *string, err error) {

	response, err := http.Get(fmt.Sprintf("%s/?q=%s&format=json", duckDuckGoAPI, url.QueryEscape(query)))
	log.Print(response)
	if err != nil {
		log.Printf("Failed to get search results from duckduckgo. %s", err)
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
func parseDuckDuckGoRecipes(html *string) (recipes *[]Recipe) {
	return &[]Recipe{}
}
