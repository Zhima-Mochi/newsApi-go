package gnews

import "time"

type News struct {
	Title           string
	Description     string
	Link            string
	Links           []string
	Content         string
	Published       string
	PublishedParsed *time.Time
	Updated         string
	UpdatedParsed   *time.Time
	GUID            string
	ImageURL        string
	Categories      []string
}
