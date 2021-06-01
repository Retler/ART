package tweets

import "sort"

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
	Sentiment float64 // TweetConsumers will enrich the Tweet structure before storing
}

type Tweets struct{
	Tweets []Tweet
}

func (t *Tweets) Sentiment() float64{
	sentiments := 0.0
	
	for _, tweet := range t.Tweets{
		sentiments += tweet.Sentiment
	}

	return sentiments / float64(len(t.Tweets))
}

type HashtagMap struct{
	Hashtag string
	Value int
}

// Returns a []HashtagMap with hashtags (keys) and the number of their occurences (values)
// The list is sorted by number of hashtag occurences
func (t *Tweets) TopXHashtags(x int) []HashtagMap{
	res := make(map[string]int)
	var hMap []HashtagMap
	
	for _, tweet := range t.Tweets{
		for _, hashtag := range tweet.Data.Entities.Hashtags{
			res[hashtag.Tag] += 1
		}
	}

	for k, v := range res{
		hMap = append(hMap, HashtagMap{k,v})
	}

	sort.Slice(hMap, func(i, j int) bool {
		return hMap[i].Value > hMap[j].Value
	})

	if len(hMap) < x{
		return hMap
	}

	return hMap[:x]
}
