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
	queryOptions = append(queryOptions, newsapi.WithLimit(10))
	endDate := time.Now()
	startDate := endDate.Add(-time.Hour * 24)
	queryOptions = append(queryOptions, newsapi.WithStartDate(startDate))
	queryOptions = append(queryOptions, newsapi.WithEndDate(endDate))
	handler.SetQueryOptions(queryOptions...)

	newsList, err := handler.GetTopNews()
	if err != nil {
		fmt.Println(err)
		return
	}
	newsapi.ToSourceLinks(newsList)
	for _, news := range newsList {
		fmt.Println("=================================")
		fmt.Println(news.Title)
		fmt.Println(news.Link)
		news.Content, err = newsapi.FetchNewsContent(news.Link)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(news.Content)
	}
}
