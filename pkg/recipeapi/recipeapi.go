package recipeapi

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	templates "masterchef_bot/pkg/helpers"

	"github.com/PuerkitoBio/goquery"
)

// Recipe object
type Recipe struct {
	Title       string
	Description string
	URL         string
}

const (
	recipeapi  = "https://duckduckgo.com"
	linkPrefix = 15
)

// SearchRecipes from duckduckgo API
func SearchRecipes(recipe string) (recipes *[]Recipe) {

	// Get search results from duckduckgo
	html, err := getDuckDuckGoSearchResult(recipe)
	if err != nil {
		return &[]Recipe{}
	}
	// Parse results to Recipes
	return parseDuckDuckGoRecipes(html)
}

// ToMessage from recipe with given title
func (recipe *Recipe) ToMessage(header string) (message string) {
	return fmt.Sprintf(templates.RecipeMessage, header, recipe.Title, recipe.URL)
}

// Get duckduckgo search result
func getDuckDuckGoSearchResult(query string) (*goquery.Document, error) {

	response, err := http.Get(fmt.Sprintf("%s/html/?q=%s+recipe", recipeapi, url.QueryEscape(query)))
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Failed to get search results from duckduckgo. %s", err)
		return nil, err
	}
	defer response.Body.Close()

	// Load the HTML document
	html, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Printf("Failed to parse search results from duckduckgo. %s", err)
		return nil, err
	}
	return html, nil
}

// Parse duckduckgo results html page
func parseDuckDuckGoRecipes(html *goquery.Document) *[]Recipe {

	var recipes []Recipe
	html.Find(".links_main.links_deep.result__body").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		url, _ := s.Find("a.result__a").Attr("href")
		title := strings.Split(s.Find("h2").Text(), " | ")[0]
		desc := s.Find(".result__snippet").Text()

		recipes = append(recipes, Recipe{
			Title:       parseWhiteSpace(title),
			Description: parseWhiteSpace(desc),
			URL:         parseURL(url),
		})
	})
	return &recipes
}

// Parse and decode URL
func parseURL(encodedURL string) string {
	decodedURL, _ := url.QueryUnescape(encodedURL)
	// Remove link prefix
	return decodedURL[linkPrefix:]
}

// Parse whitespaces
func parseWhiteSpace(in string) string {
	regex := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	return strings.TrimSpace(regex.ReplaceAllString(in, ""))
}
