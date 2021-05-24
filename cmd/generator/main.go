package main

import (
	//"fmt"
	"net/http"
	config "github.com/Retler/ART/config"
	//repo "github.com/Retler/ART/tweet_repo"
	tweets "github.com/Retler/ART/tweets"
	processing "github.com/Retler/ART/processing"
)


func main(){
	config, _ := config.GetConfig("./config/artconfig.yaml")
	tq := make(chan tweets.Tweet, 100)
	rq := make(chan processing.Result, 100)

	tp := processing.TweetProducer{
		Config: *config,
		TweetQueue: tq,
		ResultQueue: rq,
		Client: tweets.TweetClient{
			Client: http.Client{},
		},
	}

	go tp.StartStreaming()
}
