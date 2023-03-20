package main

import (
	"fmt"
	"time"

	"github.com/Zhima-Mochi/GNews-go/gnews"
)

func main() {

	google_news := gnews.NewGNews("chinese traditional", "Taiwan", 1)

	google_news.SetStartDate(time.Date(2023, 3, 19, 0, 0, 0, 0, time.UTC))
	google_news.SetEndDate(time.Date(2023, 3, 20, 0, 0, 0, 0, time.UTC))
	results, err := google_news.GetNews("technology")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, result := range results {
		fmt.Println(result.Title)
		fmt.Println(result.Link)
		_, err := gnews.GetFullArticle(result.Link)
		if err != nil {
			fmt.Println(err)
		}

	}
}
