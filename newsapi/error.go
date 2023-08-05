package newsapi

import "errors"

var (
	ErrEmptyQuery = errors.New("query cannot be empty")

	ErrEmptyTopic = errors.New("topic cannot be empty")

	ErrInvalidTopic = errors.New("invalid topic")

	ErrEmptyLocation = errors.New("location cannot be empty")

	ErrEmptyLink = errors.New("link cannot be empty")

	ErrNoSourceLink = errors.New("no source link")

	ErrFailedToGetNewsContent = errors.New("failed to get news content")
)
