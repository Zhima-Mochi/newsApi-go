package newsapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/mmcdole/gofeed"
)

var _ NewsApi = (*newsApi)(nil)

var (
	defaultNewsApi = &newsApi{
		language: "en",
		location: "US",
		limit:    10,
		client:   http.DefaultClient,
	}

	googleNewsURL = url.URL{
		Scheme: "https",
		Host:   "news.google.com",
		Path:   "/",
	}
)

type newsApi struct {
	language  string
	location  string
	period    *time.Duration
	startDate *time.Time
	endDate   *time.Time
	limit     int
	client    *http.Client
}

func NewNewsApi(options ...NewsApiOption) *newsApi {
	n := defaultNewsApi

	for _, option := range options {
		option(n)
	}

	return n
}

// SetQueryOptions sets the query options
func (n *newsApi) SetQueryOptions(options ...QueryOption) {
	for _, option := range options {
		option(n)
	}
}

// GetTopNews gets the news by path and query
func (n *newsApi) GetTopNews() ([]*News, error) {
	return n.getNews("/rss", "")
}

// GetLocationNews gets the news by location
func (n *newsApi) GetLocationNews(location string) ([]*News, error) {
	if location == "" {
		return nil, ErrEmptyLocation
	}
	path := "rss/headlines/section/geo/" + location
	return n.getNews(path, "")
}

// GetTopicNews gets the news by topic
func (n *newsApi) GetTopicNews(topic string) ([]*News, error) {
	if topic == "" {
		return nil, ErrEmptyTopic
	}
	topic = strings.ToUpper(topic)
	if _, ok := TopicMap[topic]; !ok {
		return nil, ErrInvalidTopic
	}
	path := "rss/headlines/section/topic/" + topic
	return n.getNews(path, "")
}

// SearchNews searches the news by query
func (n *newsApi) SearchNews(query string) ([]*News, error) {
	if query == "" {
		return nil, ErrEmptyQuery
	}
	query = strings.ReplaceAll(query, " ", "%20")
	return n.getNews("rss/search", query)
}

// composeURL composes the url by path and query
func (n *newsApi) composeURL(path string, query string) url.URL {
	searchURL := googleNewsURL
	q := url.Values{}
	q.Add("hl", n.language)
	q.Add("gl", n.location)
	q.Add("ceid", n.location+":"+n.language)
	searchURL.Path = path
	if query != "" {
		q.Set("q", query)
		if n.period != nil {
			q.Set("q", q.Get("q")+"+when:"+FormatDuration(*n.period))
		}
		if n.endDate != nil {
			q.Set("q", q.Get("q")+"+before:"+n.endDate.Format("2006-01-02"))
		}
		if n.startDate != nil {
			q.Set("q", q.Get("q")+"+after:"+n.startDate.Format("2006-01-02"))
		}
	}
	searchURL.RawQuery = q.Encode()
	return searchURL
}

// getNews gets the news by path and query
func (n *newsApi) getNews(path, query string) ([]*News, error) {
	searchURL := n.composeURL(path, query)
	req, err := http.NewRequest(http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", RandomUserAgent())

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	parser := gofeed.NewParser()
	feed, err := parser.ParseString(string(body))
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	newsList := make([]*News, 0, len(feed.Items))

	for _, item := range feed.Items {
		news := NewNews(item)
		newsList = append(newsList, news)
	}
	// sort by published date
	sort.Slice(newsList, func(i, j int) bool {
		return newsList[i].PublishedParsed.After(*newsList[j].PublishedParsed)
	})
	// limit the number of news
	if n.limit > 0 && n.limit < len(newsList) {
		newsList = newsList[:n.limit]
	}
	return newsList, nil
}

// FetchSourceLinks fetches the source links by the google news links
func FetchSourceLinks(newsList []*News) {
	var wg sync.WaitGroup
	for _, news := range newsList {
		wg.Add(1)
		go func(news *News) {
			defer wg.Done()
			err := news.fetchSourceLink()
			if err != nil {
				log.Println(fmt.Printf("error fetching source link: %s", err))
			}
		}(news)
	}
	wg.Wait()
}

// FetchSourceContents fetches the source contents by the source links
func FetchSourceContents(newsList []*News) {
	var wg sync.WaitGroup
	for _, news := range newsList {
		wg.Add(1)
		go func(news *News) {
			defer wg.Done()
			err := news.fetchSourceContent()
			if err != nil {
				log.Println(fmt.Printf("error fetching source content: %s", err))
			}
		}(news)
	}
	wg.Wait()
}

// Deprecated: use FetchSourceContents instead
func FetchNewsContent(link string) (string, error) {
	var content string
	if IsNewsApiLink(link) {
		var err error
		link, err = GetOriginalLink(link)
		if err != nil {
			return "", err
		}
	}
	c := colly.NewCollector(colly.Async(true))
	c.OnHTML("script", func(e *colly.HTMLElement) {
		e.DOM.Remove()
	})
	helper := func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			content += el.Text + "\n"
		})
	}
	linkURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	if selector, ok := newsHostToContentSelector[linkURL.Host]; ok {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			helper(e)
		})
	}

	err = c.Visit(link)
	if err != nil {
		return "", err
	}
	c.Wait()
	if content != "" {
		content = CleanHTML(content)
		return content, nil
	} else {
		return "", ErrFailedToGetNewsContent
	}
}

// Deprecated: use FetchSourceLinks instead
func ToSourceLinks(newsList []*News) {
	var wg sync.WaitGroup
	for _, news := range newsList {
		wg.Add(1)
		go func(news *News) {
			defer wg.Done()
			// check if the link is a google news link
			if IsNewsApiLink(news.Link) {
				originalLink, err := GetOriginalLink(news.Link)
				if err != nil {
					return
				}
				// set original link
				news.Link = originalLink
			}
		}(news)
	}
	wg.Wait()
}
