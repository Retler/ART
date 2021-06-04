package processing

import (
	"bytes"
	"errors"
	config "github.com/Retler/ART/config"
	repo "github.com/Retler/ART/tweet_repo"
	tweets "github.com/Retler/ART/tweets"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

// Tweet producer should fail in case of bad request (provoked by bad URL)
func TestTweetProducerBadUrl(t *testing.T) {
	tq := make(chan tweets.Tweet, 5)
	rq := make(chan Result, 5)
	tp := TweetProducer{
		Config:      config.Config{},
		TweetQueue:  tq,
		ResultQueue: rq,
		Client: tweets.MockClient{
			UrlFunc: func() string { return "#!¤%&" },
		},
	}

	go tp.StartStreaming()

	select {
	case tweet := <-tq:
		t.Errorf("Should not have received this tweet: %v", tweet)
	case res := <-rq:
		if res.Error == nil {
			t.Error("The result message should contain an error!")
		}
		if res.Message != "Could not create request" {
			t.Errorf("Wrong error received: %v", res.Error)
		}
	}
}

// Tweet producer should fail when request is created succesfully but fails execution
func TestTweetProducerRequestFails(t *testing.T) {
	tq := make(chan tweets.Tweet, 5)
	rq := make(chan Result, 5)
	r := &http.Response{
		StatusCode: 400,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
	}
	tp := TweetProducer{
		Config:      config.Config{},
		TweetQueue:  tq,
		ResultQueue: rq,
		Client: tweets.MockClient{
			UrlFunc: func() string { return "localhost" },
			DoFunc: func(*http.Request) (*http.Response, error) {
				return r, errors.New("Request exec error!")
			},
		},
	}

	go tp.StartStreaming()

	select {
	case tweet := <-tq:
		t.Errorf("Should not have received this tweet: %v", tweet)
	case res := <-rq:
		if res.Error == nil {
			t.Error("The result message should contain an error!")
		}
		if res.Message != "Failed to execute request" {
			t.Errorf("Wrong error received: %v", res.Error)
		}
	}
}

// Tweet producer should fail when request is executed succesfully but returns a bad status
func TestTweetProducerFailsOnBadHttpStatus(t *testing.T) {
	tq := make(chan tweets.Tweet, 5)
	rq := make(chan Result, 5)
	r := &http.Response{
		StatusCode: 400,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
	}
	tp := TweetProducer{
		Config:      config.Config{},
		TweetQueue:  tq,
		ResultQueue: rq,
		Client: tweets.MockClient{
			UrlFunc: func() string { return "localhost" },
			DoFunc: func(*http.Request) (*http.Response, error) {
				return r, nil
			},
		},
	}

	go tp.StartStreaming()

	select {
	case tweet := <-tq:
		t.Errorf("Should not have received this tweet: %v", tweet)
	case res := <-rq:
		if res.Error == nil {
			t.Error("The result message should contain an error!")
		}
		if res.Message != "Response status not OK" {
			t.Errorf("Wrong error received: %v", res.Error)
		}
	}
}

// Tweet producer should be able to read multiple Tweets from response body
func TestTweetProducerHappyPath(t *testing.T) {
	tq := make(chan tweets.Tweet, 5)
	rq := make(chan Result, 5)
	resp_body := `{"data":{"id":"1396383361833209856","lang":"en","text":"RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9","entities":{"hashtags":[{"start":1, "end":2, "tag":"BLM"}],"mentions":[{"start":3,"end":12,"username":"momy9775"}],"urls":[{"start":27,"end":50,"url":"https://t.co/6LwQHnmXK9","expanded_url":"https://twitter.com/momy9775/status/1396313215319953415/photo/1","display_url":"pic.twitter.com/6LwQHnmXK9"}]},"author_id":"1085741751174721536","created_at":"2021-05-23T08:31:51.000Z","public_metrics":{"retweet_count":12927,"reply_count":0,"like_count":0,"quote_count":0}}}` + "\n" + `{"data":{"id":"1396383361833209856","lang":"en","text":"RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9","entities":{"hashtags":[{"start":1, "end":2, "tag":"BLM"}],"mentions":[{"start":3,"end":12,"username":"momy9775"}],"urls":[{"start":27,"end":50,"url":"https://t.co/6LwQHnmXK9","expanded_url":"https://twitter.com/momy9775/status/1396313215319953415/photo/1","display_url":"pic.twitter.com/6LwQHnmXK9"}]},"author_id":"1085741751174721536","created_at":"2021-05-23T08:31:51.000Z","public_metrics":{"retweet_count":12927,"reply_count":0,"like_count":0,"quote_count":0}}}`
	expected_tweet := tweets.Tweet{
		Data: tweets.Data{
			TweetID:   "1396383361833209856",
			Content:   "RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9",
			AuthorID:  "1085741751174721536",
			CreatedAt: "2021-05-23T08:31:51.000Z",
			Language:  "en",
			PublicMetrics: tweets.PublicMetrics{
				RetweetCount: 12927,
				LikeCount:    0,
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
	r := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(resp_body))),
	}
	tp := TweetProducer{
		Config:      config.Config{},
		TweetQueue:  tq,
		ResultQueue: rq,
		Client: tweets.MockClient{
			UrlFunc: func() string { return "localhost" },
			DoFunc: func(*http.Request) (*http.Response, error) {
				return r, nil
			},
		},
	}

	go tp.StartStreaming()

	// Recieve from tweet channel twice
	for i := 0; i < 2; i++ {
		tweet := <-tq
		if !reflect.DeepEqual(expected_tweet, tweet) {
			t.Errorf("Parsed tweet didn't match. Expected:\n %+v\nGot:\n %+v\n", expected_tweet, tweet)
		}
	}

	// At last, an EOF error should be sent from tweet producer

	res, ok := <-rq
	if res.Error != io.EOF || !ok {
		t.Errorf("EOF expected, but got: %v", res.Error)
	}
}

// Test that the TweetConsumer can read tweets from the queue and store them
func TestTweetConsumer(t *testing.T) {
	tc := make(chan tweets.Tweet, 10)
	rq := make(chan Result, 10)
	tr := repo.TweetRepositoryMemory{
		Tweets: make(map[string]tweets.Tweet),
	}
	tcm := TweetConsumerSimple{
		TweetQueue:  tc,
		ResultQueue: rq,
		TweetRepo:   tr,
	}

	go tcm.StartConsuming()

	tweet := tweets.Tweet{
		Data: tweets.Data{
			TweetID:   "1396383361833209856",
			Content:   "RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9",
			AuthorID:  "1085741751174721536",
			CreatedAt: "2021-05-23T08:31:51.000Z",
			Language:  "en",
			PublicMetrics: tweets.PublicMetrics{
				RetweetCount: 12927,
				LikeCount:    0,
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

	tc <- tweet
	close(tc)

	res := <-rq
	if res.Error != nil {
		t.Error("Should err..")
	}

	consumedTweet, err := tr.GetTweet(tweet.Data.TweetID)
	if err != nil {
		t.Errorf("Could not get the tweet: %v", err)
	}
	if !reflect.DeepEqual(consumedTweet, tweet) {
		t.Errorf("Consumed tweet is not equal to the sent tweet.\nConsumed: %v\nSent: %v", consumedTweet, tweet)
	}
}

// Test that the consumer and producer can properly interact with each other on the channels
func TestConsumerAndProducer(t *testing.T) {
	resp_body := `{"data":{"id":"1396383361833209856","lang":"en","text":"RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9","entities":{"hashtags":[{"start":1, "end":2, "tag":"BLM"}],"mentions":[{"start":3,"end":12,"username":"momy9775"}],"urls":[{"start":27,"end":50,"url":"https://t.co/6LwQHnmXK9","expanded_url":"https://twitter.com/momy9775/status/1396313215319953415/photo/1","display_url":"pic.twitter.com/6LwQHnmXK9"}]},"author_id":"1085741751174721536","created_at":"2021-05-23T08:31:51.000Z","public_metrics":{"retweet_count":12927,"reply_count":0,"like_count":0,"quote_count":0}}}` + "\n" + `{"data":{"id":"1396383361833209856","lang":"en","text":"RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9","entities":{"hashtags":[{"start":1, "end":2, "tag":"BLM"}],"mentions":[{"start":3,"end":12,"username":"momy9775"}],"urls":[{"start":27,"end":50,"url":"https://t.co/6LwQHnmXK9","expanded_url":"https://twitter.com/momy9775/status/1396313215319953415/photo/1","display_url":"pic.twitter.com/6LwQHnmXK9"}]},"author_id":"1085741751174721536","created_at":"2021-05-23T08:31:51.000Z","public_metrics":{"retweet_count":12927,"reply_count":0,"like_count":0,"quote_count":0}}}`
	expected_tweet := tweets.Tweet{
		Data: tweets.Data{
			TweetID:   "1396383361833209856",
			Content:   "RT @momy9775: โควิดดีสเดย์ https://t.co/6LwQHnmXK9",
			AuthorID:  "1085741751174721536",
			CreatedAt: "2021-05-23T08:31:51.000Z",
			Language:  "en",
			PublicMetrics: tweets.PublicMetrics{
				RetweetCount: 12927,
				LikeCount:    0,
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
	r := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(resp_body))),
	}
	tq := make(chan tweets.Tweet, 5)
	rq := make(chan Result, 5)
	rq2 := make(chan Result, 5)
	tp := TweetProducer{
		Config:      config.Config{},
		TweetQueue:  tq,
		ResultQueue: rq,
		Client: tweets.MockClient{
			UrlFunc: func() string { return "localhost" },
			DoFunc: func(*http.Request) (*http.Response, error) {
				return r, nil
			},
		},
	}

	go tp.StartStreaming()

	tr := repo.TweetRepositoryMemory{
		Tweets: make(map[string]tweets.Tweet),
	}
	tcm := TweetConsumerSimple{
		TweetQueue:  tq,
		ResultQueue: rq2,
		TweetRepo:   tr,
	}

	go tcm.StartConsuming()

	_ = <-rq2 // Receive closing message (A hack to avoid listening on multiple channels)

	tweet, err := tr.GetTweet("1396383361833209856")
	if err != nil {
		t.Errorf("Could not get tweet from repository: %v", err)
	}

	if !reflect.DeepEqual(expected_tweet, tweet) {
		t.Errorf("Fetched wrong tweet. Got: %v\nExpected: %v\n", tweet, expected_tweet)
	}
}

// Test that the consumer sets the correct sentiment score (based on a prior knowledge of a score)
// This might fail if the sentiment scorer changes
func TestConsumerSentiment(t *testing.T) {
	expected_score := 0.6369499429264264
	test_text := "I love apples and coding so much"
	tweet := tweets.Tweet{
		Data: tweets.Data{
			TweetID:   "1396383361833209856",
			Content:   test_text,
			AuthorID:  "1085741751174721536",
			CreatedAt: "2021-05-23T08:31:51.000Z",
			Language:  "en",
			PublicMetrics: tweets.PublicMetrics{
				RetweetCount: 12927,
				LikeCount:    0,
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
	tq := make(chan tweets.Tweet, 5)
	rq := make(chan Result, 5)
	tr := repo.TweetRepositoryMemory{
		Tweets: make(map[string]tweets.Tweet),
	}
	tcm := TweetConsumerSimple{
		TweetQueue:  tq,
		ResultQueue: rq,
		TweetRepo:   tr,
	}

	go tcm.StartConsuming()

	tq <- tweet // Receive tweet
	close(tq)   // Close tweet channel to stop consumer
	<-rq        // Wait for result sent from consumer that marks that it is done

	tweet, err := tr.GetTweet("1396383361833209856")
	if err != nil {
		t.Errorf("Could not get tweet from repository: %v", err)
	}

	if tweet.Sentiment != expected_score {
		t.Errorf("Wrong sentiment score. Expected: %f but got %f", expected_score, tweet.Sentiment)
	}
}
