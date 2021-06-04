package main

import (
	//"fmt"
	config "github.com/Retler/ART/config"
	processing "github.com/Retler/ART/processing"
	trepo "github.com/Retler/ART/tweet_repo"
	tweets "github.com/Retler/ART/tweets"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	c, _ := config.GetConfig("./config/artconfig.yaml")
	tq := make(chan tweets.Tweet, 100)
	rq := make(chan processing.Result, 100)
	repo, err := trepo.NewMysqlRepo(*c)
	if err != nil {
		log.Fatalf("Could not initiate tweet repo: %s", err)
	}

	tp := processing.TweetProducer{
		Config:      *c,
		TweetQueue:  tq,
		ResultQueue: rq,
		Client: tweets.TweetClient{
			Client: http.Client{},
		},
	}

	tc := processing.TweetConsumerSimple{
		TweetQueue:  tq,
		ResultQueue: rq,
		TweetRepo:   repo,
	}

	go tp.StartStreaming()
	go tc.StartConsuming()

	for res := range rq {
		log.Printf("Recieved result: %v", res)
	}

}
