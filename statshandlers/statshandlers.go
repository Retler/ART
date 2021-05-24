package statshandlers

import(
	"fmt"
	"net/http"
	"encoding/json"
	repo "github.com/Retler/ART/tweet_repo"
)

func Home(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "ART StatsAPI")
}

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
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tweet)
}
