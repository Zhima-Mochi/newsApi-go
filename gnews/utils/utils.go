package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Zhima-Mochi/GNews-go/gnews/constants"
)

func LangMapping(lang string) string {
	return constants.AVAILABLE_LANGUAGES[strings.ToLower(lang)]
}

func CountryMapping(country string) string {
	return constants.AVAILABLE_COUNTRIES[strings.ToLower(country)]
}

func ProcessURL(item map[string]interface{}, excludeWebsites []string) (string, error) {
	source, ok := item["source"].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid source")
	}

	sourceHref, ok := source["href"].(string)
	if !ok {
		return "", errors.New("invalid source href")
	}

	for _, website := range excludeWebsites {
		r, _ := regexp.Compile(fmt.Sprintf(`^http(s)?://(www.)?%s.*`, strings.ToLower(website)))
		if r.MatchString(sourceHref) {
			return "", nil
		}
	}

	link, ok := item["link"].(string)
	if !ok {
		return "", errors.New("invalid link")
	}

	if matched, _ := regexp.MatchString(constants.GOOGLE_NEWS_REGEX, link); matched {
		resp, err := http.Head(link)
		if err != nil {
			return "", errors.New("error getting URL: " + err.Error())
		}
		redirectedURL, err := url.Parse(resp.Header.Get("Location"))
		if err != nil {
			return "", errors.New("error parsing URL: " + err.Error())
		}
		link = redirectedURL.String()
	}

	return link, nil
}
