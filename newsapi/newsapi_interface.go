package newsapi

type NewsApi interface {
	GetTopNews() ([]*News, error)
	GetTopicNews(topic string) ([]*News, error)
	GetLocationNews(location string) ([]*News, error)
	SearchNews(query string) ([]*News, error)

	SetQueryOptions(options ...QueryOption)
}
