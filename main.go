package main

import (
	"fmt"
	"time"

	"github.com/Zhima-Mochi/GNews-go/gnews"
)

func main() {

	google_news := gnews.NewGNews("chinese traditional", "Taiwan")
	google_news.SetLimit(1)
	before := time.Now()
	after := before.Add(-time.Hour * 24)
	google_news.SetStartDate(&after)
	google_news.SetEndDate(&before)
	_, err := google_news.GetNews()
	if err != nil {
		fmt.Println(err)
		return
	}
	// for _, result := range results {
	// _, err := gnews.GetFullArticle(result.Link)
	// if err != nil {
	// 	fmt.Println(err)
	// c := colly.NewCollector(
	// 	colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
	// )
	// c.SetRequestTimeout(10 * time.Second)
	// c.OnHTML("div.caas-body", func(e *colly.HTMLElement) {
	// 	fmt.Println(e.Text)
	// })
	// c.Visit("https://tw.news.yahoo.com/%E6%A9%9F%E5%99%A8%E4%BA%BA%E7%90%86%E8%B2%A1-%E5%8A%A9%E4%BF%A1%E8%A8%97%E8%B3%87%E7%94%A2%E9%95%B7%E5%A4%A7-201000708.html")
	// }
}
