package main

import (
	"fmt"
	"time"

	"github.com/Zhima-Mochi/GNews-go/gnews"
)

func main() {

	google_news := gnews.NewGNews("en", "US", 10)

	google_news.SetStartDate(time.Date(2023, 3, 19, 0, 0, 0, 0, time.UTC))
	google_news.SetEndDate(time.Date(2023, 3, 20, 0, 0, 0, 0, time.UTC))
	results, err := google_news.GetNews("WORLD")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, result := range results {
		fmt.Println(result.Title)
		fmt.Println(result.Description)
		fmt.Println(result.Link)
		fmt.Println(result.Published)
		fmt.Println(result.Image)
		fmt.Println()
	}
}
