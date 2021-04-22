package recipeapi

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	templates "masterchef_bot/pkg/helpers"

	"github.com/PuerkitoBio/goquery"
	"github.com/getsentry/sentry-go"
)

// Recipe object
type Recipe struct {
	Title       string
	Description string
	URL         string
}

const (
	recipeapi = "https://duckduckgo.com"
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
func getDuckDuckGoSearchResult(query string) (html *goquery.Document, err error) {

	client := &http.Client{}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/html?q=%s+recipe", recipeapi, url.QueryEscape(query)), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	response, err := client.Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		sentry.CaptureException(err)
		log.Printf("Failed to get search results from duckduckgo. %s (%v)", err, response.StatusCode)
		return nil, errors.New("Failed to get search results from duckduckgo")
	}
	defer response.Body.Close()

	// Load the HTML document
	html, err = goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Failed to parse search results from duckduckgo. %s", err)
	}
	return
}

// Parse duckduckgo results html page
func parseDuckDuckGoRecipes(html *goquery.Document) *[]Recipe {

	var recipes []Recipe
	html.Find(".links_main.links_deep.result__body").Each(func(i int, s *goquery.Selection) {
		// For each item found, get url, title and desc
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
func parseURL(encodedURL string) (parsedUrl string) {
	decodedURL, _ := url.QueryUnescape(encodedURL)
	// Remove link suffix
	parsedUrl = strings.Split(decodedURL, "&rut=")[0]
	// Remove link prefix
	regex := regexp.MustCompile(`(https).*`)
	return regex.FindString(parsedUrl)
}

// Parse whitespaces
func parseWhiteSpace(in string) string {
	regex := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	return strings.TrimSpace(regex.ReplaceAllString(in, ""))
}
