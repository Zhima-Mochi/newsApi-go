package gnews

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Zhima-Mochi/GNews-go/gnews/constants"
	"github.com/Zhima-Mochi/GNews-go/gnews/utils"
	"github.com/mmcdole/gofeed"
)

// GNews is the main struct
type GNews struct {
	Language        string
	Country         string
	MaxResults      int
	period          *time.Duration
	startDate       *time.Time
	endDate         *time.Time
	excludeWebsites *[]string
	proxy           *string
}

// NewGNews creates a new GNews instance
// Language is the language of the news (e.g. en, fr, de, etc.)
// Country is the country of the news (e.g. US, FR, DE, etc.)
// MaxResults is the maximum number of results to return
func NewGNews(Language, Country string, MaxResults int) *GNews {
	gnews := &GNews{
		Language:   Language,
		Country:    Country,
		MaxResults: MaxResults,
	}
	if language, ok := constants.AVAILABLE_LANGUAGES[Language]; ok {
		gnews.Language = language
	} else {
		gnews.Language = constants.DEFAULT_LANGUAGE
	}
	if country, ok := constants.AVAILABLE_COUNTRIES[Country]; ok {
		gnews.Country = country
	} else {
		gnews.Country = constants.DEFAULT_COUNTRY
	}
	return gnews
}

// SetPeriod sets the period of the news
// Available periods are: 1h, 1d, 7d, 30d, 1y
// If you want to set a custom period, use SetStartDate and SetEndDate
func (g *GNews) SetPeriod(period time.Duration) *GNews {
	g.period = &period
	return g
}

// SetStartDate sets the start date of the news
// If you want to set a custom period, use SetPeriod
func (g *GNews) SetStartDate(startDate time.Time) *GNews {
	g.startDate = &startDate
	return g
}

// SetEndDate sets the end date of the news
// If you want to set a custom period, use SetPeriod
func (g *GNews) SetEndDate(endDate time.Time) *GNews {
	g.endDate = &endDate
	return g
}

// SetExcludeWebsites sets the websites to exclude from the results
func (g *GNews) SetExcludeWebsites(excludeWebsites []string) *GNews {
	g.excludeWebsites = &excludeWebsites
	return g
}

// SetProxy sets the proxy to use
func (g *GNews) SetProxy(proxy string) *GNews {
	g.proxy = &proxy
	return g
}

func (g *GNews) ceid() string {
	timeQuery := ""
	if g.startDate != nil || g.endDate != nil {
		if g.period != nil {
			timeQuery += "when%3A" + g.period.String()
		}
		if g.endDate != nil {
			timeQuery += " before%3A" + g.endDate.String()
		}
		if g.startDate != nil {
			timeQuery += " after%3A" + g.startDate.String()
		}
	} else if g.period != nil {
		timeQuery += "%20when%3A" + g.period.String()
	}
	return timeQuery + "&hl=" + g.Language + "&gl=" + g.Country + "&ceid=" + g.Country + "%3A" + g.Language
}

type Article struct {
	Title       string
	Description string
	Content     string
	URL         string
	ImageURL    string
	PublishedAt time.Time
	Source      string
}

func (g *GNews) GetFullArticle(url string) (*Article, error) {
	// TODO
	return nil, nil
}

func CleanHTML(html string) string {
	// TODO
	return ""
}

func (g *GNews) process(item map[string]interface{}) (map[string]interface{}, error) {
	url, err := utils.ProcessURL(item, *g.excludeWebsites)
	if err != nil {
		return nil, err
	}

	title := item["title"].(string)

	item = map[string]interface{}{
		"title":          title,
		"description":    CleanHTML(item["description"].(string)),
		"published date": item["published"],
		"url":            url,
		"publisher":      item["source"],
	}
	return item, nil
}

// GetNews gets the news
func (g *GNews) GetNews(key string) ([]map[string]interface{}, error) {
	if key == "" {
		return nil, constants.ErrEmptyQuery
	}
	// key = ReplaceAll(key, " ", "%20")
	key = strings.ReplaceAll(key, " ", "%20")
	query := "/search?q=" + key
	return g.getNews(query)
}

func (g *GNews) getNews(query string) ([]map[string]interface{}, error) {
	link := constants.BASE_URL + query + g.ceid()
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", constants.USER_AGENT)

	if g.proxy != nil {
		proxyUrl, err := url.Parse(*g.proxy)
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy url: %w", err)
		}
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()
		fp := gofeed.NewParser()
		feed, err := fp.Parse(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error parsing feed: %w", err)
		}
		items := make([]map[string]interface{}, 0, len(feed.Items))
		for i, feedItem := range feed.Items {
			if i >= g.MaxResults {
				break
			}
			feedItemMap := map[string]interface{}{
				"title":       feedItem.Title,
				"description": feedItem.Description,
				"published":   feedItem.Published,
				"source":      feedItem.Extensions["source"]["title"][0].Value,
			}
			item, err := g.process(feedItemMap)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		return items, nil
	} else {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()
		fp := gofeed.NewParser()
		feed, err := fp.Parse(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error parsing feed: %w", err)
		}
		items := make([]map[string]interface{}, 0, len(feed.Items))
		for _, feedItem := range feed.Items {
			feedItemMap := map[string]interface{}{
				"title":       feedItem.Title,
				"description": feedItem.Description,
				"published":   feedItem.Published,
				"source":      feedItem.Extensions["source"]["title"][0].Value,
			}
			item, err := g.process(feedItemMap)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		return items, nil
	}
}

// GetTopNews gets the top news
func (g *GNews) GetTopNews() ([]map[string]interface{}, error) {
	query := "?"
	return g.getNews(query)
}

// GetNewsByTopic gets the news by topic
func (g *GNews) GetNewsByTopic(topic string) ([]map[string]interface{}, error) {
	if topic == "" {
		return nil, constants.ErrEmptyTopic
	}
	topic = strings.ToUpper(topic)
	if _, ok := constants.TOPICS[topic]; ok {
		query := "/headlines/section/topic/" + topic + "?"
		return g.getNews(query)
	} else {
		return nil, constants.ErrInvalidTopic
	}
}

// GetNewsByLocation gets the news by location
func (g *GNews) GetNewsByLocation(location string) ([]map[string]interface{}, error) {
	if location == "" {
		return nil, constants.ErrEmptyLocation
	}
	query := "/headlines/section/geo/" + location + "?"
	return g.getNews(query)
}
