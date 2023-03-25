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
> struct of [News](gnews/models.go)

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

# Example
```
Title:
德意志銀行股價重挫引發危機疑懼| 聯合新聞網 - 聯合新聞網

Link:
https://udn.com/news/story/6811/7055228

Content:
德意志銀行（Deutsche Bank）違約成本激增，重燃外界對銀行業危機擴大的疑懼，今天股價重挫。

法新社報導，德國最大銀行－德意志銀行今天在法蘭克福股市一度跌幅超過14%，收盤跌8.5%報8.54歐元。

投資人擔憂銀行業體質狀況，德意志銀行預防債務違約的風險成本、即信用違約交換（CDS）大幅攀高。

美國3家區域銀行倒閉，加上瑞士銀行集團（UBS）強制接手瑞士信貸銀行（Credit Suisse），本月稍早曾引發市場騷亂，現在德意志銀行又成為投資人關注的焦點。

德意志銀行的競爭對手－德國商業銀行（Commerzbank）也表現欠佳，今天盤中一度重挫8.5%，收盤下跌5.45%報8.88歐元。

德國銀行股走弱，領跌歐洲銀行股，法國興業銀行（Societe Generale）和法國巴黎銀行（BNP Paribas），以及英國數家銀行股價都挫跌。

德國總理蕭茲（Olaf Scholz）針對德意志銀行再度提出保證，表示這家銀行早已「現代化和組織化其營運方式，這是一家相當賺錢的銀行，沒有理由去擔心」。

蕭茲在布魯塞爾舉行的歐盟領袖峰會指出，歐洲銀行體系在嚴格規定和監管下「表現穩定」。

```

# Todo
- [ ] FetchContent() is not working properly for some news's website.
- [ ] Implement FetchAllContent(newss []*News) with goroutine.

