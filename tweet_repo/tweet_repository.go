package tweet_repo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	config "github.com/Retler/ART/config"
	tweets "github.com/Retler/ART/tweets"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
	"strings"
)

type TweetRepository interface {
	SaveTweet(tweets.Tweet) error
	GetTweet(string) (tweets.Tweet, error)
	GetTweetsSince(time.Time) (tweets.Tweets, error)
}

type TweetRepositoryMemory struct {
	Tweets map[string]tweets.Tweet
}

// Inner function which save's multiple cached tweets at once
func (trs *TweetRepositoryMysql) saveTweets(ts []tweets.Tweet) error{
	q := "INSERT INTO art.tweets VALUES "
	vals := []interface{}{}
	for _, tweet := range ts{
		hashtagBytes, err := json.Marshal(tweet.Data.Entities.Hashtags)
		hashtagString := string(hashtagBytes)
		if err != nil {
			return err
		}
		
		tweetTime, err := time.Parse(time.RFC3339, tweet.Data.CreatedAt)
		if err != nil {
			log.Errorf("Received error during parsing time: %s", err)
			return err
		}
		maxLength := math.Min(float64(len(tweet.Data.Content)), 500)
		tweetContent := tweet.Data.Content[:int(maxLength)]

		q += "(?, ?, ?, ?, ?, ?, ?, ?, ?),"
		vals = append(vals, tweet.Data.TweetID, tweet.Data.AuthorID, tweetContent, tweetTime, tweet.Data.Language, tweet.Data.PublicMetrics.RetweetCount, tweet.Data.PublicMetrics.LikeCount, hashtagString, tweet.Sentiment)
	}

	q = strings.TrimSuffix(q, ",")
	
	_, err := trs.db.Exec(q, vals...)
	if err != nil {
		log.Errorf("Received error while inserting tweets: %s", err)
		return err
	}
	
	return nil
}

func NewMysqlRepo(config config.Config, cacheSize int) (*TweetRepositoryMysql, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return &TweetRepositoryMysql{}, err
	}

	return &TweetRepositoryMysql{
		db:     db,
		config: config,
		cacheSize: cacheSize, // Nunber of tweets to batch in memory before storing
	}, nil
}

type TweetRepositoryMysql struct {
	db     *sql.DB
	config config.Config
	cache []tweets.Tweet
	cacheSize int
}

func NewMemoryRepoMock() TweetRepositoryMemory {
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

func (trs *TweetRepositoryMysql) GetTweet(tweetId string) (tweets.Tweet, error) {
	return tweets.Tweet{}, errors.New("TODO")
}

// Maps the sql 'Rows' object to the 'Tweets' model object
func (trs *TweetRepositoryMysql) rowsToTweets(rows *sql.Rows) (tweets.Tweets, error) {
	res := tweets.Tweets{}
	for rows.Next() {
		tweet := tweets.Tweet{}

		// The hashtags are stores as a json string, so we need to scan into a string first
		var hashtags string
		// The time format has to be handled explicitly
		var tweetTime string

		err := rows.Scan(&tweet.Data.TweetID, &tweet.Data.AuthorID, &tweet.Data.Content, &tweetTime, &tweet.Data.Language, &tweet.Data.PublicMetrics.RetweetCount, &tweet.Data.PublicMetrics.LikeCount, &hashtags, &tweet.Sentiment)
		if err != nil {
			return res, err
		}

		// Unmarshall hashtag json string
		err = json.Unmarshal([]byte(hashtags), &tweet.Data.Entities.Hashtags)
		if err != nil {
			return res, err
		}

		// Parse time string
		tweetTimeParsed, err := time.Parse("2006-01-02 15:04:05", tweetTime)
		if err != nil {
			return res, err
		}
		tweet.Data.CreatedAt = tweetTimeParsed.Format(time.RFC3339)

		res.Tweets = append(res.Tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (trs *TweetRepositoryMysql) GetTweetsSince(t time.Time) (tweets.Tweets, error) {
	q := `
SELECT * FROM art.tweets WHERE CREATED_AT > ?;
`
	rows, err := trs.db.Query(q, t)
	defer rows.Close()
	if err != nil {
		return tweets.Tweets{}, err
	}

	fetchedTweets, err := trs.rowsToTweets(rows)
	log.Infof("Fetched tweets: %d", len(fetchedTweets.Tweets))
	if err != nil {
		return fetchedTweets, err
	}

	return fetchedTweets, nil
}

// SaveTweet saves the tweet to the underlying mysql db. Note that the
// saving of the tweets is batched if the cacheSize > 0
func (trs *TweetRepositoryMysql) SaveTweet(tweet tweets.Tweet) error {
	trs.cache = append(trs.cache, tweet)
	log.Infof("Length of trs.cache: %d\ncacheSize: %d", len(trs.cache), trs.cacheSize)
	if len(trs.cache) > trs.cacheSize {
		err := trs.saveTweets(trs.cache)
		// Empty cache when filled. Note, all tweets are lost on error...
		trs.cache = nil 
		if err != nil{
			return err
		}
	}

	return nil
}

func (trm TweetRepositoryMemory) SaveTweet(tweet tweets.Tweet) error {
	trm.Tweets[tweet.Data.TweetID] = tweet

	return nil
}

func (trm TweetRepositoryMemory) GetTweet(tweetID string) (tweets.Tweet, error) {
	tweet, ok := trm.Tweets[tweetID]

	if !ok {
		return tweet, errors.New(fmt.Sprintf("Tweet with ID: %s could not be found", tweetID))
	}

	return tweet, nil
}

func (trm TweetRepositoryMemory) GetTweetsSince(t time.Time) (tweets.Tweets, error) {
	var res []tweets.Tweet

	for _, tweet := range trm.Tweets {
		tweetTime, err := time.Parse(time.RFC3339, tweet.Data.CreatedAt)
		if err != nil {
			return tweets.Tweets{res}, err
		}
		if tweetTime.After(t) {
			res = append(res, tweet)
		}
	}

	return tweets.Tweets{res}, nil
}
