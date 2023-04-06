package main

import (
	"fmt"
	"time"

	"github.com/Zhima-Mochi/GNews-go/gnews"
)

func main() {

	google_news := gnews.NewGNews()
	google_news.SetLimit(3)
	google_news.SetLanguage(gnews.LanguageChineseTraditional)
	google_news.SetLocation(gnews.LocationTaiwan)
	before := time.Now()
	after := before.Add(-time.Hour * 24)
	google_news.SetStartDate(&after)
	google_news.SetEndDate(&before)
	newss, err := google_news.GetTopNews()
	if err != nil {
		fmt.Println(err)
		return
	}
	google_news.ConvertToOriginalLinks(newss)
	for _, news := range newss {
		fmt.Println("=================================")
		fmt.Println(news.Title)
		fmt.Println(news.Link)
		_, err := news.FetchContent()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(news.Content)
	}
}
