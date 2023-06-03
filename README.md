# Google News API
[![Go Report Card](https://goreportcard.com/badge/github.com/Zhima-Mochi/newsApi-go)](https://goreportcard.com/report/github.com/Zhima-Mochi/newsApi-go)

The News API is a Go package that allows you to fetch news articles from Google News. It provides a simple and convenient way to retrieve news based on various criteria such as language, location, topic, and search query.

## Installation

To use the News API package in your Go project, you can install it using the `go get` command:

```
go get github.com/Zhima-Mochi/newsApi-go/newsapi
```

## Usage

Import the News API package in your Go code:

```go
import "github.com/Zhima-Mochi/newsApi-go/newsapi"
```

### Creating a News API instance

You can create a new instance of the News API by calling the `NewNewsApi` function. You can also provide optional configuration options to customize the behavior of the API.

```go
api := newsapi.NewNewsApi()
```

### Fetching top news

To retrieve the top news articles, you can use the `GetTopNews` method:

```go
newsList, err := api.GetTopNews()
if err != nil {
    // handle error
}

// Process the news articles
for _, news := range newsList {
    // Access news properties such as title, description, link, etc.
    fmt.Println(news.Title)
}
```

### Fetching news by location

You can retrieve news articles based on a specific location using the `GetLocationNews` method:

```go
newsList, err := api.GetLocationNews(newsapi.LocationUnitedStates)
if err != nil {
    // handle error
}

// Process the news articles
for _, news := range newsList {
    // Access news properties
    fmt.Println(news.Title)
}
```

### Fetching news by topic

To fetch news articles related to a specific topic, you can use the `GetTopicNews` method:

```go
newsList, err := api.GetTopicNews(newsapi.TopicTechnology)
if err != nil {
    // handle error
}

// Process the news articles
for _, news := range newsList {
    // Access news properties
    fmt.Println(news.Title)
}
```

### Searching for news

You can search for news articles using a specific query using the `SearchNews` method:

```go
newsList, err := api.SearchNews("Go programming language")
if err != nil {
    // handle error
}

// Process the news articles
for _, news := range newsList {
    // Access news properties
    fmt.Println(news.Title)
}
```

### Customizing the API options

The News API provides various options to customize the behavior of the API. You can set the query options using the `SetQueryOptions` method:

```go
api.SetQueryOptions(
    newsapi.WithLanguage("en"),
    newsapi.WithLocation("US"),
    newsapi.WithLimit(20),
)

// Fetch news based on the configured options
newsList, err := api.GetTopNews()
// ...
```

### Fetching content of a news article

```go
newsContent, err := newsapi.FetchNewsContent(news.Link)
if err != nil {
    // handle error
}

// Access news content
fmt.Println(newsContent)
```



## Example

Please refer to the [example](example/main.go) for a complete example of using the News API package.

## Todo
- [ ] FetchNewsContent() is not working properly for some news's website.
- [ ] Implement FetchAllNewsContent(newsList []*News) with goroutine.

## License

The News API package is open source and available under the [MIT License](https://github.com/Zhima-Mochi/newsApi-go/blob/main/LICENSE).