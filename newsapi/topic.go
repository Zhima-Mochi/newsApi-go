package newsapi

const (
	// Topic
	TopicWorld         string = "WORLD"
	TopicNation        string = "NATION"
	TopicBusiness      string = "BUSINESS"
	TopicTechnology    string = "TECHNOLOGY"
	TopicEntertainment string = "ENTERTAINMENT"
	TopicSports        string = "SPORTS"
	TopicScience       string = "SCIENCE"
	TopicHealth        string = "HEALTH"
)

var (
	TopicMap = map[string]string{
		TopicWorld:         "w",
		TopicNation:        "n",
		TopicBusiness:      "b",
		TopicTechnology:    "t",
		TopicEntertainment: "e",
		TopicSports:        "s",
		TopicScience:       "snc",
		TopicHealth:        "m",
	}
)
