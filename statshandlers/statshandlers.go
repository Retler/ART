package statshandlers

import (
	"encoding/json"
	"fmt"
	repo "github.com/Retler/ART/tweet_repo"
	tweets "github.com/Retler/ART/tweets"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ART StatsAPI")
}

func GetTweet(w http.ResponseWriter, r *http.Request) {
	tweets := repo.NewMemoryRepoMock()
	keys, ok := r.URL.Query()["tweet_id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(400)
		return
	}

	tweetID := keys[0]
	tweet, err := tweets.GetTweet(tweetID)
	if err != nil {
		log.Errorf("Got error: %s", err)
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tweet)
}

func GetTweetsSince(w http.ResponseWriter, r *http.Request) {
	tweets := repo.NewMemoryRepoMock()
	keys, ok := r.URL.Query()["age_minutes"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(400)
		return
	}

	ageMinutesStr := keys[0]
	ageMinutes, err := strconv.Atoi(ageMinutesStr)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	ageOfTweets := time.Now().UTC().Add(-time.Duration(ageMinutes) * time.Minute)
	tweetsSince, err := tweets.GetTweetsSince(ageOfTweets)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if len(tweetsSince.Tweets) == 0 {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tweetsSince.Tweets)
}

type StatRes struct {
	Sentiment   map[string]float64
	TopHashtags []tweets.HashtagMap
}

type Time struct {
	Amount   int
	Duration time.Duration
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	tweets := repo.NewMemoryRepoMock()
	unit, okUnit := r.URL.Query()["unit"]

	if !okUnit || (unit[0] != "minute" && unit[0] != "second" && unit[0] != "hour") {
		w.WriteHeader(400)
		return
	}

	unitStr := unit[0]
	ageUnitMap := map[string]Time{
		"second": Time{60, time.Second},
		"minute": Time{60, time.Minute},
		"hour":   Time{10000, time.Hour},
	}

	ageOfTweets := time.Now().UTC().Add(-time.Duration(ageUnitMap[unitStr].Amount) * ageUnitMap[unitStr].Duration)
	tweetsSince, err := tweets.GetTweetsSince(ageOfTweets)
	if err != nil {
		log.Errorf("Got error: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(tweetsSince.Tweets) == 0 {
		w.WriteHeader(404)
		return
	}

	bins, err := tweetsSince.ByDuration(ageUnitMap[unitStr].Duration)
	if err != nil {
		log.Errorf("Got error: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := StatRes{
		Sentiment:   bins.Sentiment(),
		TopHashtags: tweetsSince.TopXHashtags(10),
	}

	json.NewEncoder(w).Encode(res)
}
