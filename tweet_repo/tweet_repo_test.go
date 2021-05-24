package tweet_repo

import(
	"time"
	"reflect"
	"testing"
	tweets "github.com/Retler/ART/tweets"
)

func TestMemoryRepo(t *testing.T){
	repo := TweetRepositoryMemory{
		Tweets: make(map[string]tweets.Tweet),
	}

	tweet := tweets.Tweet{
		Data: tweets.Data{
			TweetID: "1396383361833209856",
			Content: "RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9",
			AuthorID: "1085741751174721536",
			CreatedAt: "2021-05-23T08:31:51.000Z",
			Language: "en",
			PublicMetrics: tweets.PublicMetrics{
				RetweetCount: 12927,
				LikeCount: 0,
			},
			Entities: tweets.Entities{
				Hashtags: []tweets.Hashtag{
					tweets.Hashtag{
						Tag: "BLM",
					},
				},
			},
		},
	}

	tweet2 := tweets.Tweet{
		Data: tweets.Data{
			TweetID: "1396383361833209857",
			Content: "RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9",
			AuthorID: "1085741751174721536",
			CreatedAt: "2021-04-23T08:31:51.000Z",
			Language: "en",
			PublicMetrics: tweets.PublicMetrics{
				RetweetCount: 12927,
				LikeCount: 0,
			},
			Entities: tweets.Entities{
				Hashtags: []tweets.Hashtag{
					tweets.Hashtag{
						Tag: "BLM",
					},
				},
			},
		},
	}

	tweet3 := tweets.Tweet{
		Data: tweets.Data{
			TweetID: "1396383361833209858",
			Content: "RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9",
			AuthorID: "1085741751174721536",
			CreatedAt: "2021-03-23T08:31:51.000Z",
			Language: "en",
			PublicMetrics: tweets.PublicMetrics{
				RetweetCount: 12927,
				LikeCount: 0,
			},
			Entities: tweets.Entities{
				Hashtags: []tweets.Hashtag{
					tweets.Hashtag{
						Tag: "BLM",
					},
				},
			},
		},
	}

	err := repo.SaveTweet(tweet)
	if err != nil{
		t.Errorf("Error on saving tweet: %v", err)
	}

	fetchedTweet, err := repo.GetTweet(tweet.Data.TweetID)
	if err != nil{
		t.Errorf("Error getting tweet: %v", err)
	}

	if !reflect.DeepEqual(tweet, fetchedTweet){
		t.Errorf("Expected tweet %v but got %v", tweet, fetchedTweet)
	}

	repo.SaveTweet(tweet2)
	repo.SaveTweet(tweet3)

	test_time, err := time.Parse(time.RFC3339, "2021-03-24T08:31:51.000Z")
	if err != nil{
		t.Errorf("Could not parse time: %v", err)
	}
	
	tweets, err := repo.GetTweetsSince(test_time)
	if err != nil{
		t.Errorf("Got error during 'GetTweetsSince': %v", err)
	}

	if len(tweets) != 2 {
		t.Errorf("Expected 2 tweets but got: %d", len(tweets))
	}

}
