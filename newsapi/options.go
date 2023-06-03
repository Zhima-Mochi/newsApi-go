package newsapi

import (
	"net/http"
	"net/url"
	"time"
)

type QueryOption func(*newsApi)

func WithLanguage(language string) QueryOption {
	return func(n *newsApi) {
		n.language = language
	}
}

func WithLocation(location string) QueryOption {
	return func(n *newsApi) {
		n.location = location
	}
}

func WithLimit(limit int) QueryOption {
	if limit > MaxSearchResults {
		limit = MaxSearchResults
	}
	return func(n *newsApi) {
		n.limit = limit
	}
}

func WithPeriod(period time.Duration) QueryOption {
	return func(n *newsApi) {
		n.period = &period
	}
}

func WithoutPeriod() QueryOption {
	return func(n *newsApi) {
		n.period = nil
	}
}

func WithStartDate(startDate time.Time) QueryOption {
	return func(n *newsApi) {
		n.startDate = &startDate
	}
}

func WithoutStartDate() QueryOption {
	return func(n *newsApi) {
		n.startDate = nil
	}
}

func WithEndDate(endDate time.Time) QueryOption {
	return func(n *newsApi) {
		n.endDate = &endDate
	}
}

func WithoutEndDate() QueryOption {
	return func(n *newsApi) {
		n.endDate = nil
	}
}

func WithoutDuration() QueryOption {
	return func(n *newsApi) {
		n.period = nil
		n.startDate = nil
		n.endDate = nil
	}
}

type NewsApiOption func(*newsApi)

func WithProxy(proxy *url.URL) NewsApiOption {
	return func(n *newsApi) {
		n.client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}
}

func WithoutProxy() NewsApiOption {
	return func(n *newsApi) {
		n.client = nil
	}
}
