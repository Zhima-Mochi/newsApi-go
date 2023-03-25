package gnews

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	collector       *colly.Collector
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
	collector := colly.NewCollector(colly.Async(true))
	gnews := &GNews{
		baseURL:   baseURL,
		language:  language,
		country:   country,
		limit:     utils.MaxSearchResults,
		collector: collector,
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

func (g *GNews) composeURL(search string) url.URL {
	searchURL := g.baseURL
	query := url.Values{}
	query.Add("hl", g.language)
	query.Add("gl", g.country)
	query.Add("ceid", g.country+":"+g.language)
	if search == "" {
		searchURL.Path = "rss"
	} else {
		searchURL.Path = "rss/search"
		query.Set("q", search)
		if g.period != nil {
			query.Set("q", query.Get("q")+"+when:"+g.period.String())
		}
		if g.endDate != nil {
			query.Set("q", query.Get("q")+"+before:"+g.endDate.Format("2006-01-02"))
		}
		if g.startDate != nil {
			query.Set("q", query.Get("q")+"+after:"+g.startDate.Format("2006-01-02"))
		}
	}
	searchURL.RawQuery = query.Encode()
	return searchURL
}

// type Article struct {
// 	Title       string
// 	Description string
// 	Content     string
// 	URL         string
// 	ImageURL    string
// 	PublishedAt time.Time
// 	Source      string
// }

// func GetFullArticle(url string) (*Article, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, errors.New("failed to download article")
// 	}
// 	defer resp.Body.Close()

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		return nil, errors.New("failed to parse HTML")
// 	}
// 	text, err := html2text.FromReader(strings.NewReader(doc.Text()), html2text.Options{})
// 	if err != nil {
// 		return nil, errors.New("failed to convert HTML to text")
// 	}
// 	article := &Article{
// 		Title:       doc.Find("h1").Text(),
// 		Description: doc.Find("meta[name=description]").AttrOr("content", ""),
// 		Content:     text,
// 		URL:         url,
// 		ImageURL:    doc.Find("meta[property='og:image']").AttrOr("content", ""),
// 		PublishedAt: time.Now(),
// 		Source:      doc.Find("meta[property='og:site_name']").AttrOr("content", ""),
// 	}

// 	return article, nil
// }

func (g *GNews) cleanRSSItem(item *gofeed.Item) (*gofeed.Item, error) {
	item.Description = utils.CleanDescription(item.Description)
	return item, nil
}

// GetNewsWithSearch gets the news with a search query
func (g *GNews) GetNewsWithSearch(search string) ([]*gofeed.Item, error) {
	if search == "" {
		return nil, utils.ErrEmptyQuery
	}
	search = strings.ReplaceAll(search, " ", "%20")
	return g.getNews(search)
}

// GetNews gets the news
func (g *GNews) GetNews() ([]*gofeed.Item, error) {
	return g.getNews("")
}

// get original news link from google news
func (g *GNews) getOriginalLink(sourceLink string) (string, error) {
	originalLink := ""
	g.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		originalLink = e.Attr("href")
	})
	err := g.collector.Visit(sourceLink)
	if err != nil {
		return "", err
	}
	g.collector.Wait()
	return originalLink, nil
}

func (g *GNews) GetItems(client *http.Client, req *http.Request) ([]*gofeed.Item, error) {
	feedItems, err := utils.GetFeedItems(client, req)
	if err != nil {
		return nil, err
	}
	items := make([]*gofeed.Item, 0, len(feedItems))
	itemCount := 0
	for _, feedItem := range feedItems {
		if itemCount >= g.limit {
			break
		}
		originalLink, err := g.getOriginalLink(feedItem.Link)
		if err != nil {
			return nil, err
		}
		if utils.IsExcludedSource(originalLink, g.excludeWebsites) {
			continue
		}
		// set original link
		feedItem.Link = originalLink
		fmt.Println(feedItem.Link)
		item, err := g.cleanRSSItem(feedItem)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		itemCount++
	}
	return items, nil
}

func (g *GNews) getNews(search string) ([]*gofeed.Item, error) {
	searchURL := g.composeURL(search)
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
		return nil, utils.ErrEmptyTopic
	}
	topic = strings.ToUpper(topic)
	if _, ok := utils.TOPICS[topic]; ok {
		query := "/headlines/section/topic/" + topic + "?"
		return g.getNews(query)
	} else {
		return nil, utils.ErrInvalidTopic
	}
}

// GetNewsByLocation gets the news by location
func (g *GNews) GetNewsByLocation(location string) ([]*gofeed.Item, error) {
	if location == "" {
		return nil, utils.ErrEmptyLocation
	}
	query := "/headlines/section/geo/" + location + "?"
	return g.getNews(query)
}
