package statshandlers

import(
	"fmt"
	"time"
	"strconv"
	"net/http"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	repo "github.com/Retler/ART/tweet_repo"
	tweets "github.com/Retler/ART/tweets"
)

func Home(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "ART StatsAPI")
}

// TODO: Change API endpoints to serve stats rather than tweets
func GetTweet(w http.ResponseWriter, r *http.Request){
	tweets := repo.NewMemoryRepoMock()
	keys, ok := r.URL.Query()["tweet_id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(400)
		return
	}

	tweetID := keys[0]
	tweet, err := tweets.GetTweet(tweetID)
	if err != nil{
		log.Errorf("Got error: %s", err)
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tweet)
}

func GetTweetsSince(w http.ResponseWriter, r *http.Request){
	tweets := repo.NewMemoryRepoMock()
	keys, ok := r.URL.Query()["age_minutes"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(400)
		return
	}

	ageMinutesStr := keys[0]
	ageMinutes, err := strconv.Atoi(ageMinutesStr)
	if err != nil{
		w.WriteHeader(404)
		return
	}
	
	ageOfTweets := time.Now().UTC().Add(-time.Duration(ageMinutes) * time.Minute)
	tweetsSince, err := tweets.GetTweetsSince(ageOfTweets)
	if err != nil{
		w.WriteHeader(500)
		return
	}

	if len(tweetsSince.Tweets) == 0{
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tweetsSince.Tweets)
}

type StatRes struct{
	Sentiment map[string]string
	TopHashtags []tweets.HashtagMap
}

func GetStats(w http.ResponseWriter, r *http.Request){
	tweets := repo.NewMemoryRepoMock()
	keys, ok := r.URL.Query()["age_minutes"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(400)
		return
	}

	ageMinutesStr := keys[0]
	ageMinutes, err := strconv.Atoi(ageMinutesStr)
	if err != nil{
		log.Errorf("Got error: %s", err)
		w.WriteHeader(404)
		return
	}
	
	ageOfTweets := time.Now().UTC().Add(-time.Duration(ageMinutes) * time.Minute)
	tweetsSince, err := tweets.GetTweetsSince(ageOfTweets)
	if err != nil{
		log.Errorf("Got error: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(tweetsSince.Tweets) == 0{
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := StatRes{
		Sentiment: map[string]string{
			"sentiment": fmt.Sprintf("%f", tweetsSince.Sentiment()),
			"n_tweets": fmt.Sprintf("%d", len(tweetsSince.Tweets)),
		},
		TopHashtags: tweetsSince.TopXHashtags(10),
	}
	
	json.NewEncoder(w).Encode(res)
}
