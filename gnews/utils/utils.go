package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Zhima-Mochi/GNews-go/gnews/constants"
	"github.com/mmcdole/gofeed"
)

func LangMapping(lang string) string {
	return constants.AVAILABLE_LANGUAGES[strings.ToLower(lang)]
}

func CountryMapping(country string) string {
	return constants.AVAILABLE_COUNTRIES[strings.ToLower(country)]
}

func ProcessURL(item *gofeed.Item, excludeWebsites *[]string) (string, error) {

	if excludeWebsites != nil {
		for _, website := range *excludeWebsites {
			r, _ := regexp.Compile(fmt.Sprintf(`^http(s)?://(www.)?%s.*`, strings.ToLower(website)))
			if r.MatchString(item.Link) {
				return "", nil
			}
		}
	}

	// Check if the item.Link is a Google News link
	link, err := url.Parse(item.Link)
	if err != nil {
		return "", errors.New("error parsing URL: " + err.Error())
	}
	// If the item.Link is a Google News link, get the real link
	if matched, _ := regexp.MatchString(constants.GOOGLE_NEWS_REGEX, link.Host); matched {
		resp, err := http.Head(link.Host)
		if err != nil {
			return "", errors.New("error getting URL: " + err.Error())
		}
		redirectedURL, err := url.Parse(resp.Header.Get("Location"))
		if err != nil {
			return "", errors.New("error parsing URL: " + err.Error())
		}
		link = redirectedURL
	}

	return link.String(), nil
}
