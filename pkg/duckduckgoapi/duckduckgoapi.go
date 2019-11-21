package duckduckgoapi

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// Recipe fron duckduckgo API
type Recipe struct {
	Title        string
	Description  string
	ThumbnailURL string
	URL          string
}

const (
	duckDuckGoAPI = "https://duckduckgo.com"
)

// SearchRecipes from duckduckgo API
func SearchRecipes(recipe string) (recipes *[]Recipe) {

	// Get search ressults from duckduckgo
	html, err := getDuckDuckGoSearchResult(recipe)
	if err != nil {
		return &[]Recipe{}
	}
	// Parse results to Recipes
	return parseDuckDuckGoRecipes(html)
}

// Get duckduckgo search result
func getDuckDuckGoSearchResult(query string) (*goquery.Document, error) {

	response, err := http.Get(fmt.Sprintf("%s/html/?q=%s+recipe&ia=recipes&iax=recipes", duckDuckGoAPI, url.QueryEscape(query)))
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Failed to get search results from duckduckgo. %s", err)
		return nil, err
	}
	defer response.Body.Close()

	// Load the HTML document
	html, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Printf("Failed to parse search results from duckduckgo. %s", err)
	}
	return html, nil
}

// Parse google results html page
func parseDuckDuckGoRecipes(html *goquery.Document) (recipes *[]Recipe) {
	// Find the review items
	html.Find(".result__body").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		url, _ := s.Find("a").Attr("href")
		title := s.Find("h2").Text()
		log.Printf("Review %d: %s - %s\n", i, url, title)
	})

	return &[]Recipe{}
}
