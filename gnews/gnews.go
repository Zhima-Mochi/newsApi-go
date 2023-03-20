package gnews

import (
	"fmt"
	"io/ioutil"
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
			timeQuery += "%20when%3A" + g.period.String()
		}
		if g.endDate != nil {
			timeQuery += "%20before%3A" + g.endDate.Format("2006-01-02")
		}
		if g.startDate != nil {
			timeQuery += "%20after%3A" + g.startDate.Format("2006-01-02")
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

func (g *GNews) process(item *gofeed.Item) (*gofeed.Item, error) {
	url, err := utils.ProcessURL(item, g.excludeWebsites)
	if err != nil {
		return nil, err
	}

	item.Link = url
	item.Description = CleanHTML(item.Description)

	return item, nil
}

// GetNews gets the news
func (g *GNews) GetNews(key string) ([]*gofeed.Item, error) {
	if key == "" {
		return nil, constants.ErrEmptyQuery
	}
	key = strings.ReplaceAll(key, " ", "%20")
	query := "/search?q=" + key
	return g.getNews(query)
}

func (g *GNews) GetItems(client *http.Client, req *http.Request) ([]*gofeed.Item, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(string(body))
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}
	items := make([]*gofeed.Item, 0, len(feed.Items))
	for _, feedItem := range feed.Items {
		item, err := g.process(feedItem)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (g *GNews) getNews(query string) ([]*gofeed.Item, error) {
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
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		client := &http.Client{
			Transport: transport,
		}
		items, err := g.GetItems(client, req)
		if err != nil {
			return nil, err
		}
		return items, nil
	} else {
		items, err := g.GetItems(http.DefaultClient, req)
		if err != nil {
			return nil, err
		}
		return items, nil
	}
}

// GetTopNews gets the top news
func (g *GNews) GetTopNews() ([]*gofeed.Item, error) {
	query := "?"
	return g.getNews(query)
}

// GetNewsByTopic gets the news by topic
func (g *GNews) GetNewsByTopic(topic string) ([]*gofeed.Item, error) {
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
func (g *GNews) GetNewsByLocation(location string) ([]*gofeed.Item, error) {
	if location == "" {
		return nil, constants.ErrEmptyLocation
	}
	query := "/headlines/section/geo/" + location + "?"
	return g.getNews(query)
}
