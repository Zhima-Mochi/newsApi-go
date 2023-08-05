package newsapi

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/mmcdole/gofeed"
)

var ()

type News struct {
	Title           string
	Description     string
	Link            string
	Links           []string
	Content         string
	Published       string
	PublishedParsed *time.Time
	Updated         string
	UpdatedParsed   *time.Time
	GUID            string
	ImageURL        string
	Categories      []string

	SourceLink        string
	SourceTitle       string
	SourceImageURL    string
	SourceImageWidth  int
	SourceImageHeight int
	SourceDescription string
	SourceKeywords    []string
	SourceSiteName    string
	SourceContent     string
}

func NewNews(item *gofeed.Item) *News {
	n := &News{
		Title:           item.Title,
		Description:     item.Description,
		Link:            item.Link,
		Links:           item.Links,
		Published:       item.Published,
		PublishedParsed: item.PublishedParsed,
		Updated:         item.Updated,
		UpdatedParsed:   item.UpdatedParsed,
		GUID:            item.GUID,
		Categories:      item.Categories,
	}

	if item.Image != nil {
		n.ImageURL = item.Image.URL
	}
	n.Description = CleanHTML(n.Description)
	return n
}

func (n *News) fetchSourceLink() error {
	if n.SourceLink != "" {
		return nil
	}

	// check if the link is a google news link
	if IsNewsApiLink(n.Link) {
		originalLink, err := GetOriginalLink(n.Link)
		if err != nil {
			return fmt.Errorf("error getting original link: %w", err)
		}
		// set source link
		n.SourceLink = originalLink
	}
	return nil
}

func (n *News) fetchSourceContent() error {
	if n.SourceContent != "" {
		return nil
	}

	if n.SourceLink == "" {
		err := n.fetchSourceLink()
		if err != nil {
			return fmt.Errorf("error fetching source link: %s", err)
		}
	}
	if n.SourceLink == "" {
		return ErrNoSourceLink
	}

	var content string
	c := colly.NewCollector(colly.Async(true))
	// remove script tag
	c.OnHTML("script", func(e *colly.HTMLElement) {
		e.DOM.Remove()
	})

	helper := func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			content += el.Text + "\n"
		})
	}

	linkURL, err := url.Parse(n.SourceLink)
	if err != nil {
		return fmt.Errorf("error parsing source link: %w", err)
	}

	// Set a callback for when a meta tag with the property "og:title" is encountered
	c.OnHTML(`meta[property="og:title"]`, func(e *colly.HTMLElement) {
		ogTitle := e.Attr("content")
		n.SourceTitle = ogTitle
	})

	// Set a callback for when a meta tag with the property "og:image" is encountered
	c.OnHTML(`meta[property="og:image"]`, func(e *colly.HTMLElement) {
		ogImage := e.Attr("content")
		n.SourceImageURL = ogImage
	})

	// Set a callback for when a meta tag with the property "og:image:width" is encountered
	c.OnHTML(`meta[property="og:image:width"]`, func(e *colly.HTMLElement) {
		ogImageWidth := e.Attr("content")
		i, _ := strconv.Atoi(ogImageWidth)
		n.SourceImageWidth = i
	})

	// Set a callback for when a meta tag with the property "og:image:height" is encountered
	c.OnHTML(`meta[property="og:image:height"]`, func(e *colly.HTMLElement) {
		ogImageHeight := e.Attr("content")
		i, _ := strconv.Atoi(ogImageHeight)
		n.SourceImageHeight = i
	})

	// Set a callback for when a meta tag with the property "og:description" is encountered
	c.OnHTML(`meta[property="og:description"]`, func(e *colly.HTMLElement) {
		ogDescription := e.Attr("content")
		n.SourceDescription = ogDescription
	})

	// Set a callback for when a meta tag with the property "og:site_name" is encountered
	c.OnHTML(`meta[property="og:site_name"]`, func(e *colly.HTMLElement) {
		ogSiteName := e.Attr("content")
		n.SourceSiteName = ogSiteName
	})

	// Set a callback for when a meta tag with the property "og:keywords" is encountered
	c.OnHTML(`meta[property="og:keywords"]`, func(e *colly.HTMLElement) {
		ogKeywords := e.Attr("content")
		n.SourceKeywords = strings.Split(ogKeywords, ",")
	})

	if selector, ok := newsHostToContentSelector[linkURL.Host]; ok {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			helper(e)
		})
	}

	// visit the source link
	err = c.Visit(n.SourceLink)
	if err != nil {
		return fmt.Errorf("error visiting source link: %w", err)
	}
	c.Wait()
	if content != "" {
		content = CleanHTML(content)
		n.SourceContent = content
	} else {
		return ErrFailedToGetNewsContent
	}
	return nil
}
