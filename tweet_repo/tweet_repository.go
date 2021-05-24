package tweet_repo

import (
	"fmt"
	"errors"
	"time"
	tweets "github.com/Retler/ART/tweets"
)

type TweetRepository interface{
	SaveTweet(tweets.Tweet) error
	GetTweet(string) (tweets.Tweet, error)
	GetTweetsSince(time.Time) ([]tweets.Tweet, error)
}

type TweetRepositoryMemory struct{
	Tweets map[string]tweets.Tweet
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

func(trm TweetRepositoryMemory) GetTweetsSince(t time.Time) ([]tweets.Tweet, error){
	var res []tweets.Tweet

	for _, tweet := range trm.Tweets {
		tweetTime, err := time.Parse(time.RFC3339, tweet.Data.CreatedAt)
		if err != nil{
			return res, err
		}
		if tweetTime.After(t){
			res = append(res, tweet)
		}
	}

	return res, nil
}
