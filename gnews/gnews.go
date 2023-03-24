package gnews

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Zhima-Mochi/GNews-go/gnews/utils"
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

func CleanHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	return doc.Text()
}

func (g *GNews) process(item *gofeed.Item) (*gofeed.Item, error) {
	url, err := utils.ProcessURL(item, g.excludeWebsites)
	if err != nil {
		return nil, err
	}
	fmt.Println(item.Description)
	item.Link = url
	item.Description = CleanHTML(item.Description)
	item.Content = CleanHTML(item.Content)
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
	for i, feedItem := range feed.Items {
		if i >= g.limit {
			break
		}
		item, err := g.process(feedItem)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	return items, nil
}

func (g *GNews) getNews(search string) ([]*gofeed.Item, error) {
	searchURL := g.composeURL(search)
	req, err := http.NewRequest(http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", utils.USER_AGENT)
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
