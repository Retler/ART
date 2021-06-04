package tweets

import (
	"fmt"
	"sort"
	"time"
)

func (ttb TweetTimeBins) Sentiment() map[string]float64 {
	res := make(map[string]float64)
	for k, bin := range ttb {
		tweets := Tweets{bin}
		res[k] = tweets.Sentiment()
	}

	return res
}

func (t *Tweets) Sentiment() float64 {
	sentiments := 0.0

	for _, tweet := range t.Tweets {
		sentiments += tweet.Sentiment
	}

	return sentiments / float64(len(t.Tweets))
}

func (t *Tweets) ByDuration(d time.Duration) (TweetTimeBins, error) {
	bins := make(TweetTimeBins)
	for _, tweet := range t.Tweets {
		tweetTime, err := time.Parse(time.RFC3339, tweet.Data.CreatedAt)
		if err != nil {
			return nil, err
		}

		label := fmt.Sprintf("%d", tweetTime.Truncate(d).Unix())

		bins[label] = append(bins[label], tweet)
	}

	return bins, nil
}

// Returns a []HashtagMap with hashtags (keys) and the number of their occurences (values)
// The list is sorted by number of hashtag occurences
func (t *Tweets) TopXHashtags(x int) []HashtagMap {
	res := make(map[string]int)
	var hMap []HashtagMap

	for _, tweet := range t.Tweets {
		for _, hashtag := range tweet.Data.Entities.Hashtags {
			res[hashtag.Tag] += 1
		}
	}

	for k, v := range res {
		hMap = append(hMap, HashtagMap{k, v})
	}

	sort.Slice(hMap, func(i, j int) bool {
		return hMap[i].Value > hMap[j].Value
	})

	if len(hMap) < x {
		return hMap
	}

	return hMap[:x]
}
