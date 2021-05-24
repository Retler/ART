package tweets

type Hashtag struct{
	Tag string `json:"tag"`
}

type Entities struct{
	Hashtags []Hashtag `json:"hashtags"`
}

type PublicMetrics struct{
	RetweetCount int `json:"retweet_count"`
	LikeCount int `json:"like_count"`
}

type Data struct{
	TweetID string `json:"id"`// The ID of the tweet
	Language string `json:"lang"` // The detected language of the tweet
	Entities Entities `json:"entities"` // The Twitter API 'entities' object (used to get hashtags)
	Content string `json:"text"`// The contents of the tweet
	AuthorID string `json:"author_id"` // The author of the Tweet
	CreatedAt string `json:"created_at"` // ISO 8601 formatted date
	PublicMetrics PublicMetrics `json:"public_metrics"` // The Twitter API `public_metrics` object
}

type Tweet struct{
	Data Data `json:"data"`
}
