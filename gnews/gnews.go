package gnews

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Zhima-Mochi/GNews-go/gnews/utils"
	"github.com/gocolly/colly"
)

const (
	// Topic
	TopicWorld         string = "WORLD"
	TopicNation        string = "NATION"
	TopicBusiness      string = "BUSINESS"
	TopicTechnology    string = "TECHNOLOGY"
	TopicEntertainment string = "ENTERTAINMENT"
	TopicSports        string = "SPORTS"
	TopicScience       string = "SCIENCE"
	TopicHealth        string = "HEALTH"
)

var (
	TopicMap = map[string]string{
		TopicWorld:         "w",
		TopicNation:        "n",
		TopicBusiness:      "b",
		TopicTechnology:    "t",
		TopicEntertainment: "e",
		TopicSports:        "s",
		TopicScience:       "snc",
		TopicHealth:        "m",
	}
)

// Host Language
const (
	LanguageEnglish            = "en"
	LanguageIndonesian         = "id"
	LanguageCzech              = "cs"
	LanguageGerman             = "de"
	LanguageSpanish            = "es-419"
	LanguageFrench             = "fr"
	LanguageItalian            = "it"
	LanguageLatvian            = "lv"
	LanguageLithuanian         = "lt"
	LanguageHungarian          = "hu"
	LanguageDutch              = "nl"
	LanguageNorwegian          = "no"
	LanguagePolish             = "pl"
	LanguagePortugueseBrasil   = "pt-419"
	LanguagePortuguesePortugal = "pt-150"
	LanguageRomanian           = "ro"
	LanguageSlovak             = "sk"
	LanguageSlovenian          = "sl"
	LanguageSwedish            = "sv"
	LanguageVietnamese         = "vi"
	LanguageTurkish            = "tr"
	LanguageGreek              = "el"
	LanguageBulgarian          = "bg"
	LanguageRussian            = "ru"
	LanguageSerbian            = "sr"
	LanguageUkrainian          = "uk"
	LanguageHebrew             = "he"
	LanguageArabic             = "ar"
	LanguageMarathi            = "mr"
	LanguageHindi              = "hi"
	LanguageBengali            = "bn"
	LanguageTamil              = "ta"
	LanguageTelugu             = "te"
	LanguageMalyalam           = "ml"
	LanguageThai               = "th"
	LanguageChineseSimplified  = "zh-Hans"
	LanguageChineseTraditional = "zh-Hant"
	LanguageJapanese           = "ja"
	LanguageKorean             = "ko"
)

// Geographic Location
const (
	LocationAustralia          = "AU"
	LocationBotswana           = "BW"
	LocationCanada             = "CA"
	LocationEthiopia           = "ET"
	LocationGhana              = "GH"
	LocationIndia              = "IN"
	LocationIndonesia          = "ID"
	LocationIreland            = "IE"
	LocationIsrael             = "IL"
	LocationKenya              = "KE"
	LocationLatvia             = "LV"
	LocationMalaysia           = "MY"
	LocationNamibia            = "NA"
	LocationNewZealand         = "NZ"
	LocationNigeria            = "NG"
	LocationPakistan           = "PK"
	LocationPhilippines        = "PH"
	LocationSingapore          = "SG"
	LocationSouthAfrica        = "ZA"
	LocationTanzania           = "TZ"
	LocationUganda             = "UG"
	LocationUnitedKingdom      = "GB"
	LocationUnitedStates       = "US"
	LocationZimbabwe           = "ZW"
	LocationCzechRepublic      = "CZ"
	LocationGermany            = "DE"
	LocationAustria            = "AT"
	LocationSwitzerland        = "CH"
	LocationArgentina          = "AR"
	LocationChile              = "CL"
	LocationColombia           = "CO"
	LocationCuba               = "CU"
	LocationMexico             = "MX"
	LocationPeru               = "PE"
	LocationVenezuela          = "VE"
	LocationBelgium            = "BE"
	LocationFrance             = "FR"
	LocationMorocco            = "MA"
	LocationSenegal            = "SN"
	LocationItaly              = "IT"
	LocationLithuania          = "LT"
	LocationHungary            = "HU"
	LocationNetherlands        = "NL"
	LocationNorway             = "NO"
	LocationPoland             = "PL"
	LocationBrazil             = "BR"
	LocationPortugal           = "PT"
	LocationRomania            = "RO"
	LocationSlovakia           = "SK"
	LocationSlovenia           = "SI"
	LocationSweden             = "SE"
	LocationVietnam            = "VN"
	LocationTurkey             = "TR"
	LocationGreece             = "GR"
	LocationBulgaria           = "BG"
	LocationRussia             = "RU"
	LocationUkraine            = "UA"
	LocationSerbia             = "RS"
	LocationUnitedArabEmirates = "AE"
	LocationSaudiArabia        = "SA"
	LocationLebanon            = "LB"
	LocationEgypt              = "EG"
	LocationBangladesh         = "BD"
	LocationThailand           = "TH"
	LocationChina              = "CN"
	LocationTaiwan             = "TW"
	LocationHongKong           = "HK"
	LocationJapan              = "JP"
	LocationRepublicOfKorea    = "KR"
)

// GNews is the main struct
type GNews struct {
	baseURL         url.URL
	language        string
	location        string
	period          *time.Duration
	startDate       *time.Time
	endDate         *time.Time
	excludeWebsites *[]string
	proxy           *string
	limit           int
}

// NewGNews creates a new GNews instance
// Language and location are optional
// If you don't specify them, the default language and location will be used
// The default language is traditional Chinese
// The default location is TW
// Language is the language of the news (e.g. en, fr, de, etc.)
// Location is the location of the news (e.g. US, FR, DE, etc.)
func NewGNews() *GNews {
	baseURL := url.URL{
		Scheme: "https",
		Host:   "news.google.com",
		Path:   "/",
	}
	gnews := &GNews{
		baseURL:  baseURL,
		language: LanguageChineseTraditional,
		location: LocationTaiwan,
		limit:    utils.MaxSearchResults,
	}
	return gnews
}

// SetLanguage sets the language of the news
func (g *GNews) SetLanguage(language string) *GNews {
	g.language = language
	return g
}

// SetLocation sets the location of the news
func (g *GNews) SetLocation(location string) *GNews {
	g.location = location
	return g
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

func (n *News) FetchContent() (string, error) {
	var content string
	url := n.Link
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
	} else if strings.Contains(url, "udn.com") {
		c.OnHTML(".article-content__editor", func(e *colly.HTMLElement) {
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
	} else if strings.Contains(url, "cnn.com") {
		c.OnHTML(".zn-body__paragraph", func(e *colly.HTMLElement) {
			helper(e)
		})
	} else if strings.Contains(url, "reuters.com") {
		c.OnHTML(".StandardArticleBody_body", func(e *colly.HTMLElement) {
			helper(e)
		})
	} else if strings.Contains(url, "cnbc.com") {
		c.OnHTML(".group", func(e *colly.HTMLElement) {
			helper(e)
		})
	} else if strings.Contains(url, "marketwatch.com") {
		c.OnHTML(".article__body", func(e *colly.HTMLElement) {
			helper(e)
		})
	} else {
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
		n.Content = content
		return n.Content, nil
	} else {
		return "", utils.ErrFailedToGetNewsContent
	}
}

// GetTop gets the top news
func (g *GNews) GetTopNews() ([]*News, error) {
	return g.getNews("rss", "")
}

// Search searches the news
func (g *GNews) SearchNews(query string) ([]*News, error) {
	if query == "" {
		return nil, utils.ErrEmptyQuery
	}
	query = strings.ReplaceAll(query, " ", "%20")
	return g.getNews("rss/search", query)
}

// GetLocationNews gets the news by location
func (g *GNews) GetLocationNews(location string) ([]*News, error) {
	if location == "" {
		return nil, utils.ErrEmptyLocation
	}
	path := "rss/headlines/section/geo/" + location
	return g.getNews(path, "")
}

// GetNewsByTopic gets the news by topic
func (g *GNews) GetTopicNews(topic string) ([]*News, error) {
	if topic == "" {
		return nil, utils.ErrEmptyTopic
	}
	topic = strings.ToUpper(topic)
	if _, ok := TopicMap[topic]; ok {
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
	q.Add("gl", g.location)
	q.Add("ceid", g.location+":"+g.language)
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

func (g *GNews) getItems(client *http.Client, req *http.Request) ([]*News, error) {
	feedItems, err := utils.GetFeedItems(client, req)
	if err != nil {
		return nil, err
	}
	items := make([]*News, 0, len(feedItems))
	itemCount := 0
	var wg sync.WaitGroup
	for _, feedItem := range feedItems {
		if itemCount >= g.limit {
			break
		}
		news := &News{
			Title:           feedItem.Title,
			Description:     feedItem.Description,
			Link:            feedItem.Link,
			Links:           feedItem.Links,
			Published:       feedItem.Published,
			PublishedParsed: feedItem.PublishedParsed,
			Updated:         feedItem.Updated,
			UpdatedParsed:   feedItem.UpdatedParsed,
			GUID:            feedItem.GUID,
			Categories:      feedItem.Categories,
		}
		if feedItem.Image != nil {
			news.ImageURL = feedItem.Image.URL
		}
		wg.Add(1)
		go func(news *News) {
			defer wg.Done()
			originalLink, err := g.getOriginalLink(news.Link)
			if err != nil {
				return
			}
			if utils.IsExcludedSource(originalLink, g.excludeWebsites) {
				return
			}
			// set original link
			news.Link = originalLink
			item, err := g.cleanRSSItem(news)
			if err != nil {
				return
			}
			items = append(items, item)
			itemCount++
		}(news)
	}
	wg.Wait()
	return items, nil
}

func (g *GNews) getNews(path, query string) ([]*News, error) {
	searchURL := g.composeURL(path, query)
	req, err := http.NewRequest(http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", utils.RandomUserAgent())
	items := make([]*News, 0)
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
		items, err = g.getItems(client, req)
		if err != nil {
			return nil, err
		}
		return items, nil
	} else {
		items, err = g.getItems(http.DefaultClient, req)
		if err != nil {
			return nil, err
		}
	}
	// sort by published date
	sort.Slice(items, func(i, j int) bool {
		return items[i].PublishedParsed.After(*items[j].PublishedParsed)
	})
	return items, nil
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

func (g *GNews) cleanRSSItem(item *News) (*News, error) {
	item.Description = utils.CleanDescription(item.Description)
	return item, nil
}
