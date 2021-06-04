package tweets

import (
	"encoding/json"
	"reflect"
	"testing"
)

// Test that a json object from the Twitter streaming API can be parsed into the local struct
func TestTweetParsing(t *testing.T) {
	stream_output := `{"data":{"id":"1396383361833209856","lang":"en","text":"RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9","entities":{"hashtags":[{"start":1, "end":2, "tag":"BLM"}],"mentions":[{"start":3,"end":12,"username":"momy9775"}],"urls":[{"start":27,"end":50,"url":"https://t.co/6LwQHnmXK9","expanded_url":"https://twitter.com/momy9775/status/1396313215319953415/photo/1","display_url":"pic.twitter.com/6LwQHnmXK9"}]},"author_id":"1085741751174721536","created_at":"2021-05-23T08:31:51.000Z","public_metrics":{"retweet_count":12927,"reply_count":0,"like_count":0,"quote_count":0}}}`

	expected_tweet := Tweet{
		Data: Data{
			TweetID:   "1396383361833209856",
			Content:   "RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9",
			AuthorID:  "1085741751174721536",
			CreatedAt: "2021-05-23T08:31:51.000Z",
			Language:  "en",
			PublicMetrics: PublicMetrics{
				RetweetCount: 12927,
				LikeCount:    0,
			},
			Entities: Entities{
				Hashtags: []Hashtag{
					Hashtag{
						Tag: "BLM",
					},
				},
			},
		},
	}

	tweet := Tweet{}

	err := json.Unmarshal([]byte(stream_output), &tweet)
	if err != nil {
		t.Errorf("Could not unmarshal stream output. Error: %v ", err)
	}

	if !reflect.DeepEqual(expected_tweet, tweet) {
		t.Errorf("Parsed tweet didn't match. Expected:\n %+v\nGot:\n %+v\n", expected_tweet, tweet)
	}
}
