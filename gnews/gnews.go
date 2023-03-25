package gnews

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Zhima-Mochi/GNews-go/gnews/utils"
	"github.com/gocolly/colly"
	"github.com/mmcdole/gofeed"
)

// GNews is the main struct
type GNews struct {
	baseURL         url.URL
	language        string
	country         string
	period          *time.Duration
	startDate       *time.Time
	endDate         *time.Time
	excludeWebsites *[]string
	proxy           *string
	limit           int
}

// NewGNews creates a new GNews instance
// Language and country are optional
// If you don't specify them, the default language and country will be used
// The default language is traditional Chinese
// The default country is TW
// Language is the language of the news (e.g. en, fr, de, etc.)
// Country is the country of the news (e.g. US, FR, DE, etc.)
func NewGNews(language string, country string) *GNews {
	baseURL := url.URL{
		Scheme: "https",
		Host:   "news.google.com",
		Path:   "/",
	}
	if lang, ok := utils.AVAILABLE_LANGUAGES[language]; ok {
		language = lang
	} else {
		language = utils.DEFAULT_LANGUAGE
	}
	if ctry, ok := utils.AVAILABLE_COUNTRIES[country]; ok {
		country = ctry
	} else {
		country = utils.DEFAULT_COUNTRY
	}
	gnews := &GNews{
		baseURL:  baseURL,
		language: language,
		country:  country,
		limit:    utils.MaxSearchResults,
	}
	return gnews
}

// SetLimit sets the limit of the results
func (g *GNews) SetLimit(limit int) *GNews {
	if limit > utils.MaxSearchResults {
		limit = utils.MaxSearchResults
	}
	g.limit = limit
	return g
}

// SetPeriod sets the period of the news
// Available periods are: 1h, 1d, 7d, 30d, 1y
// If you want to set a custom period, use SetStartDate and SetEndDate
func (g *GNews) SetPeriod(period *time.Duration) *GNews {
	g.period = period
	return g
}

// SetStartDate sets the start date of the news
// If you want to set a custom period, use SetPeriod
func (g *GNews) SetStartDate(startDate *time.Time) *GNews {
	g.startDate = startDate
	return g
}

// SetEndDate sets the end date of the news
// If you want to set a custom period, use SetPeriod
func (g *GNews) SetEndDate(endDate *time.Time) *GNews {
	g.endDate = endDate
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

func (g *GNews) GetNewsContent(url string) (string, error) {
	var content string
	c := colly.NewCollector(colly.Async(true))
	c.OnHTML("script", func(e *colly.HTMLElement) {
		e.DOM.Remove()
	})
	helper := func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			content += el.Text + "\n"
		})
	}

	if strings.Contains(url, "yahoo.com") {
		c.OnHTML(".caas-body", func(e *colly.HTMLElement) {
			helper(e)
		})
	} else if strings.Contains(url, "chinatimes.com") {
		c.OnHTML(".article-body", func(e *colly.HTMLElement) {
			helper(e)
		})
	} else if strings.Contains(url, "tvbs.com") {
		c.OnHTML(".article_content", func(e *colly.HTMLElement) {
			helper(e)
		})
	} else if strings.Contains(url, "appledaily.com") {
		c.OnHTML(".ndArticle_margin", func(e *colly.HTMLElement) {
			helper(e)
		})
	} else if strings.Contains(url, "ettoday.net") {
		c.OnHTML(".story", func(e *colly.HTMLElement) {
			helper(e)
		})
	} else if strings.Contains(url, "ltn.com") {
		c.OnHTML(".text", func(e *colly.HTMLElement) {
			helper(e)
		})
	} else if strings.Contains(url, "chinatimes.com") {
		c.OnHTML(".article-body", func(e *colly.HTMLElement) {
			helper(e)
		})
	}

	err := c.Visit(url)
	if err != nil {
		return "", err
	}
	c.Wait()
	if content != "" {
		content = utils.CleanHTML(content)
		return content, nil
	} else {
		return "", utils.ErrFailedToGetNewsContent
	}
}

// GetTopNews gets the top news
func (g *GNews) GetTopNews() ([]*gofeed.Item, error) {
	return g.getNews("rss", "")
}

// GetNewsWithSearch gets the news with search
func (g *GNews) GetNewsWithSearch(query string) ([]*gofeed.Item, error) {
	if query == "" {
		return nil, utils.ErrEmptyQuery
	}
	query = strings.ReplaceAll(query, " ", "%20")
	return g.getNews("rss/search", query)
}

// GetNewsByLocation gets the news by location
func (g *GNews) GetNewsByLocation(location string) ([]*gofeed.Item, error) {
	if location == "" {
		return nil, utils.ErrEmptyLocation
	}
	path := "rss/headlines/section/geo/" + location
	return g.getNews(path, "")
}

// GetNewsByTopic gets the news by topic
func (g *GNews) GetNewsByTopic(topic string) ([]*gofeed.Item, error) {
	if topic == "" {
		return nil, utils.ErrEmptyTopic
	}
	topic = strings.ToUpper(topic)
	if _, ok := utils.TOPICS[topic]; ok {
		path := "rss/headlines/section/topic/" + topic
		return g.getNews(path, "")
	} else {
		return nil, utils.ErrInvalidTopic
	}
}

func (g *GNews) composeURL(path, query string) url.URL {
	searchURL := g.baseURL
	q := url.Values{}
	q.Add("hl", g.language)
	q.Add("gl", g.country)
	q.Add("ceid", g.country+":"+g.language)
	searchURL.Path = path
	if query != "" {
		q.Set("q", query)
		if g.period != nil {
			q.Set("q", q.Get("q")+"+when:"+g.period.String())
		}
		if g.endDate != nil {
			q.Set("q", q.Get("q")+"+before:"+g.endDate.Format("2006-01-02"))
		}
		if g.startDate != nil {
			q.Set("q", q.Get("q")+"+after:"+g.startDate.Format("2006-01-02"))
		}
	}
	searchURL.RawQuery = q.Encode()
	return searchURL
}

func (g *GNews) getItems(client *http.Client, req *http.Request) ([]*gofeed.Item, error) {
	feedItems, err := utils.GetFeedItems(client, req)
	if err != nil {
		return nil, err
	}
	items := make([]*gofeed.Item, 0, len(feedItems))
	itemCount := 0
	var wg sync.WaitGroup
	for _, feedItem := range feedItems {
		if itemCount >= g.limit {
			break
		}
		wg.Add(1)
		go func(feedItem *gofeed.Item) {
			defer wg.Done()
			originalLink, err := g.getOriginalLink(feedItem.Link)
			if err != nil {
				return
			}
			if utils.IsExcludedSource(originalLink, g.excludeWebsites) {
				return
			}
			// set original link
			feedItem.Link = originalLink
			item, err := g.cleanRSSItem(feedItem)
			if err != nil {
				return
			}
			items = append(items, item)
			itemCount++
		}(feedItem)
	}
	wg.Wait()
	return items, nil
}

func (g *GNews) getNews(path, query string) ([]*gofeed.Item, error) {
	searchURL := g.composeURL(path, query)
	req, err := http.NewRequest(http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", utils.RandomUserAgent())
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
		items, err := g.getItems(client, req)
		if err != nil {
			return nil, err
		}
		return items, nil
	} else {
		items, err := g.getItems(http.DefaultClient, req)
		if err != nil {
			return nil, err
		}
		return items, nil
	}
}

// get original news link from google news
func (g *GNews) getOriginalLink(sourceLink string) (string, error) {
	originalLink := ""
	c := colly.NewCollector(colly.Async(true))
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		originalLink = e.Attr("href")
	})
	err := c.Visit(sourceLink)
	if err != nil {
		return "", err
	}
	c.Wait()
	return originalLink, nil
}

func (g *GNews) cleanRSSItem(item *gofeed.Item) (*gofeed.Item, error) {
	item.Description = utils.CleanDescription(item.Description)
	return item, nil
}
