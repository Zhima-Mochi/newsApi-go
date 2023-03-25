package main

import (
	"fmt"
	"time"

	"github.com/Zhima-Mochi/GNews-go/gnews"
)

func main() {

	google_news := gnews.NewGNews()
	google_news.SetLimit(10)
	google_news.SetLanguage(gnews.LanguageEnglish)
	google_news.SetLocation(gnews.LocationTaiwan)
	before := time.Now()
	after := before.Add(-time.Hour * 24)
	google_news.SetStartDate(&after)
	google_news.SetEndDate(&before)
	results, err := google_news.GetTopicNews(gnews.TopicBusiness)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, result := range results {
		// content, err := google_news.GetNewsContent(result.Link)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// fmt.Println(content)
		fmt.Println(result.Title)
		fmt.Println(result.Link)
		fmt.Println(result.Published)
		fmt.Println(result.Updated)
		fmt.Println(result.Extensions)
		fmt.Println("  ")
		fmt.Println("  ")
	}
}
