package smartraveller

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pariz/gountries"
)

const smartravellerAllAdvisoriesRssFeedURL = "https://www.smartraveller.gov.au/countries/documents/index.rss"

// :(
var customCountryPageSlugToAlpha2 = map[string]string{
	"brunei-darussalam":                           "BN",
	"cote-divoire-ivory-coast":                    "CI",
	"democratic-republic-congo":                   "CD",
	"eswatini":                                    "SZ",
	"federated-states-micronesia":                 "FM",
	"israel-and-occupied-palestinian-territories": "PS", // IL
	"kosovo": "XK",
	"macau":  "MO",
	"north-korea-democratic-peoples-republic-korea": "KP",
	"north-macedonia":            "MK",
	"south-korea-republic-korea": "KR",
	"timor-leste":                "TL",
	"turkiye":                    "TR",
	"united-states-america":      "US",
}

type Country struct {
	Name   string `json:"name"`   // The name of the country
	Alpha2 string `json:"alpha2"` // The ISO Alpha 2 code of the country
}

type Advisory struct {
	Country      Country `json:"country"`      // The country object containing name and alpha2 code
	Advice       string  `json:"advice"`       // The advisory text
	LatestUpdate string  `json:"latestUpdate"` // The update text

	Level     int    `json:"level"`     // The advisory level (1-5)
	Published string `json:"published"` // The published date of the advisory
	PageUrl   string `json:"pageUrl"`   // The URL of the advisory page
}

func getUrlSlug(url string) string {
	// Split the URL by slashes
	parts := strings.Split(url, "/")
	// Return the last part of the URL
	return parts[len(parts)-1]
}

func removeTextInBrackets(country string) string {
	re := regexp.MustCompile(`\s*\(.*?\)`)
	return re.ReplaceAllString(country, "")
}

func cleanString(str string) string {
	// Removes HTML tags and new lines from the string
	re := regexp.MustCompile(`<[^>]*>`)
	str = re.ReplaceAllString(str, "")

	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, "&nbsp;", " ")

	re = regexp.MustCompile(`\s+`)
	str = re.ReplaceAllString(str, " ")

	return str
}

func parseDate(date string) string {
	// The date looks like "27 Mar 2025 23:00:00 AEDT" convert to UTC RFC3339
	// This isn't ideal but the Smartraveller RSS contains abbreviated TZ names
	// The standard timezone database does not perform mapping from abbrviation to a zone offset.

	// Replace the timezone abbreviations with their respective offsets
	date = strings.Replace(date, "AEDT", "+11:00", 1)
	date = strings.Replace(date, "AEST", "+10:00", 1)
	date = strings.Replace(date, "ACDT", "+10:30", 1)
	date = strings.Replace(date, "ACST", "+09:30", 1)

	layout := "02 Jan 2006 15:04:05 -07:00"
	parsedTime, err := time.Parse(layout, date)
	if err != nil {
		// Handle parse error
		fmt.Printf("Error parsing date: %v\n", err)
		return ""
	}

	return parsedTime.UTC().Format(time.RFC3339) // Convert to UTC and format as RFC3339
}

func parseCountry(country string, pageUrl string) (Country, error) {
	// Frustratingly the Smartraveller API does not provide the country name in a standard format
	// The country will look like "China", "Spain", we need to convert it to a standard format (ISO Alpha 2)
	// Use various methods to try map the country name

	country = removeTextInBrackets(country)
	query := gountries.New()

	// Attempt to find the country by name
	countryData, err := query.FindCountryByName(country)
	if err == nil {
		return Country{
			Name:   countryData.Name.Common,
			Alpha2: countryData.Alpha2,
		}, nil
	}

	// If not found, try using the URL slug to map to alpha2 code
	urlSlug := getUrlSlug(pageUrl)
	if alpha2, exists := customCountryPageSlugToAlpha2[urlSlug]; exists {
		countryData, err = query.FindCountryByAlpha(alpha2)
		if err == nil {
			return Country{
				Name:   countryData.Name.Common,
				Alpha2: countryData.Alpha2,
			}, nil
		}
	}

	// As a fallback, replace dashes in the slug with spaces and try finding by name again
	urlSlug = strings.ReplaceAll(urlSlug, "-", " ")
	countryData, err = query.FindCountryByName(urlSlug)
	if err != nil {
		return Country{}, err
	}

	return Country{
		Name:   countryData.Name.Common,
		Alpha2: countryData.Alpha2,
	}, nil
}

func parseAdvisoryLevel(levelValue string) (int, error) {
	// The levelValue will look like "5/5"
	if len(levelValue) >= 3 {
		levelValue = levelValue[0:1]
		level, err := strconv.Atoi(levelValue)
		if err != nil {
			return 0, err
		} else if level < 1 || level > 5 {
			return 0, err
		}
		return level, nil
	}
	return 0, errors.New("could not parse advisory level")
}

func parseAdvisory(item gofeed.Item) (Advisory, error) {
	// The actual advisory levels are nested in a custom "ta:warnings" object
	stTaWarnings := item.Extensions["ta"]["warnings"][0].Children

	levelValue := stTaWarnings["level"][0].Value
	advisoryValue := stTaWarnings["description"][0].Value

	level, levelErr := parseAdvisoryLevel(levelValue)
	country, countryErr := parseCountry(item.Title, item.Link)

	if levelErr != nil || countryErr != nil {
		log.Printf("Error parsing advisory: %v, %v", levelErr, countryErr)
		return Advisory{}, errors.New("could not parse advisory")
	}

	advisory := Advisory{
		Country:      country,
		Advice:       advisoryValue,
		LatestUpdate: cleanString(item.Description),
		Level:        level,

		Published: parseDate(item.Published),
		PageUrl:   item.Link,
	}

	return advisory, nil
}

var ErrAdvisoryNotFound = errors.New("advisory not found")

func GetAdvisories(country string) ([]Advisory, error) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(smartravellerAllAdvisoriesRssFeedURL)

	if country != "" {
		country = strings.ToUpper(country)
	}

	var advisories []Advisory

	for _, item := range feed.Items {
		advisory, err := parseAdvisory(*item)
		if err != nil {
			continue
		}

		if country != "" && advisory.Country.Alpha2 == country {
			return []Advisory{advisory}, nil
		}

		advisories = append(advisories, advisory)
	}

	if country != "" {
		return nil, ErrAdvisoryNotFound
	}

	return advisories, nil
}
