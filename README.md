# GNews-go

[![Go Report Card](https://goreportcard.com/badge/github.com/Zhima-Mochi/Gnews-go)](https://goreportcard.com/report/github.com/Zhima-Mochi/Gnews-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


# Introduction
This repository contains a Go implementation of the [ranahaani/GNews](https://github.com/ranahaani/GNews) code.

# Installation
```bash
go get github.com/Zhima-Mochi/Gnews-go
```

# Usage

```go
package main

import (
    "fmt"
    "github.com/Zhima-Mochi/Gnews-go"
)

func main(){
    gnews := gnews.NewGnews()
    // return struct of News
    newss, err := gnews.GetTopNews() 
}
```

## Set options
```go
gnews.SetLocation(gnews.CountryJapan) // default is Taiwan
gnews.SetLanguage(gnews.LanguageJapanese) // default is Traditional Chinese
gnews.SetBefore(time.Now())
gnews.SetAfter(time.Now().AddDate(0, 0, -7))
gnews.SetMaxResults(10) // default is 100
```

## Search by keyword
```go
newss, err := gnews.SearchNews("keyword")
```
## Fetch content
We do not have the content of the news in the original data, so we need to fetch it from the news website in the `Link` field.
```go
content, err := news.FetchContent()
```
