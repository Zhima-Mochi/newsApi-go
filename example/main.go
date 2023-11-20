package main

import (
	"fmt"
	"time"

	"github.com/Zhima-Mochi/newsApi-go/newsapi"
)

func main() {
	handler := newsapi.NewNewsApi()
	queryOptions := []newsapi.QueryOption{}
	queryOptions = append(queryOptions, newsapi.WithLanguage(newsapi.LanguageChineseTraditional))
	queryOptions = append(queryOptions, newsapi.WithLocation(newsapi.LocationTaiwan))
	queryOptions = append(queryOptions, newsapi.WithLimit(3))
	endDate := time.Now()
	startDate := endDate.Add(-time.Hour * 72)
	queryOptions = append(queryOptions, newsapi.WithStartDate(startDate))
	queryOptions = append(queryOptions, newsapi.WithEndDate(endDate))
	// queryOptions = append(queryOptions, newsapi.WithPeriod(time.Hour))
	handler.SetQueryOptions(queryOptions...)

	newsList, err := handler.GetTopNews()
	if err != nil {
		fmt.Println(err)
		return
	}
	newsapi.FetchSourceLinks(newsList)
	for _, news := range newsList {
		fmt.Println("=================================")
		fmt.Println(news.Title)
		fmt.Println(news.Link)
		// news.Content, err = newsapi.FetchNewsContent(news.Link)
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }
		newsapi.FetchSourceContents([]*newsapi.News{news})
		// fmt.Println(news.Content)
		fmt.Println(news.SourceLink)
		fmt.Println(news.SourceTitle)
		fmt.Println(news.SourceImageURL)
		fmt.Println(news.SourceImageWidth)
		fmt.Println(news.SourceImageHeight)
		fmt.Println(news.SourceDescription)
		fmt.Println(news.SourceKeywords)
		fmt.Println(news.SourceSiteName)
		fmt.Println(news.SourceIconUrl)
		fmt.Println(news.SourceContent)
	}
}
