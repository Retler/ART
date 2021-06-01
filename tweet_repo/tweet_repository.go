package tweet_repo

import (
	"fmt"
	"errors"
	"time"
	log "github.com/sirupsen/logrus"
	tweets "github.com/Retler/ART/tweets"
)

type TweetRepository interface{
	SaveTweet(tweets.Tweet) error
	GetTweet(string) (tweets.Tweet, error)
	GetTweetsSince(time.Time) (tweets.Tweets, error)
}

type TweetRepositoryMemory struct{
	Tweets map[string]tweets.Tweet
}

func NewMemoryRepoMock() TweetRepositoryMemory{
	tweets := map[string]tweets.Tweet{
		tweets.MockTweet1.Data.TweetID: tweets.MockTweet1,
		tweets.MockTweet2.Data.TweetID: tweets.MockTweet2,
		tweets.MockTweet3.Data.TweetID: tweets.MockTweet3,
		tweets.MockTweet4.Data.TweetID: tweets.MockTweet4,
	}
	
	return TweetRepositoryMemory{
		Tweets: tweets,
	}
}


func(trm TweetRepositoryMemory) SaveTweet(tweet tweets.Tweet) error{
	trm.Tweets[tweet.Data.TweetID] = tweet

	return nil
}

func(trm TweetRepositoryMemory) GetTweet(tweetID string) (tweets.Tweet, error){
	tweet, ok := trm.Tweets[tweetID] 

	if !ok{
		return tweet, errors.New(fmt.Sprintf("Tweet with ID: %s could not be found", tweetID))
	}
	
	return tweet, nil
}

func(trm TweetRepositoryMemory) GetTweetsSince(t time.Time) (tweets.Tweets, error){
	var res []tweets.Tweet

	for _, tweet := range trm.Tweets {
		tweetTime, err := time.Parse(time.RFC3339, tweet.Data.CreatedAt)
		if err != nil{
			return tweets.Tweets{res}, err
		}
		if tweetTime.After(t){
			res = append(res, tweet)
		}
	}

	return tweets.Tweets{res}, nil
}
