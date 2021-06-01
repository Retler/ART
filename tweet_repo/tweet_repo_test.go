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

	err := repo.SaveTweet(tweets.MockTweet1)
	if err != nil{
		t.Errorf("Error on saving tweet: %v", err)
	}

	fetchedTweet, err := repo.GetTweet(tweets.MockTweet1.Data.TweetID)
	if err != nil{
		t.Errorf("Error getting tweet: %v", err)
	}

	if !reflect.DeepEqual(tweets.MockTweet1, fetchedTweet){
		t.Errorf("Expected tweet %v but got %v", tweets.MockTweet1, fetchedTweet)
	}

	repo.SaveTweet(tweets.MockTweet2)
	repo.SaveTweet(tweets.MockTweet3)

	test_time, err := time.Parse(time.RFC3339, "2021-03-24T08:31:51.000Z")
	if err != nil{
		t.Errorf("Could not parse time: %v", err)
	}
	
	fetchedTweets, err := repo.GetTweetsSince(test_time)
	if err != nil{
		t.Errorf("Got error during 'GetTweetsSince': %v", err)
	}

	if len(fetchedTweets.Tweets) != 2 {
		t.Errorf("Expected 2 tweets but got: %d", len(fetchedTweets.Tweets))
	}

}
