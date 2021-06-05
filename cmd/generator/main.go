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
	rq1 := make(chan processing.Result, 100)
	rq2 := make(chan processing.Result, 100)
	repo, err := trepo.NewMysqlRepo(*c)
	if err != nil {
		log.Fatalf("Could not initiate tweet repo: %s", err)
	}

	tp := processing.TweetProducer{
		Config:      *c,
		TweetQueue:  tq,
		ResultQueue: rq1,
		Client: tweets.TweetClient{
			Client: http.Client{},
		},
	}

	tc := processing.TweetConsumerSimple{
		TweetQueue:  tq,
		ResultQueue: rq2,
		TweetRepo:   repo,
	}

	go tp.StartStreaming()
	go tc.StartConsuming()

	for {
		select {
		case res1, ok := <-rq1:
			log.Printf("Recieved producer result: %v", res1)
			if !ok {
				log.Info("rq1 closed. 'nilling' channel")
				rq1 = nil // Stop this channel from being selected
			}
		case res2, ok := <-rq2:
			log.Printf("Recieved consumer result: %v", res2)
			if !ok {
				log.Info("rq2 closed. 'nilling' channel")
				rq2 = nil // Stop this channel from being selected
			}
		}

		if rq1 == nil && rq2 == nil {
			break // Both channels are closed. Save to exit now.
		}
	}

}
