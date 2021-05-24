# ART (Analyzer of Realtime Tweets)
ART is a system which provides real-time statistics on Tweets consumed from the [1% Twitter sampling API](https://developer.twitter.com/en/docs/twitter-api/tweets/sampled-stream/api-reference/get-tweets-sample-stream).

The repository is a work-in-progress. At the moment, only the Twitter integration is implemented (Tweet Producer, Tweet Consumer and an in-memory Tweet Repository). The final architecture will consist of:

* **Tweet Producer:** Parses and queues Tweets from the sample stream
* **Tweet Consumer:** Processes tweets (enriching with additional information) and stores them
* **Tweet Repository:** A place where consumers can store the tweets
* **Stats API:** An API connected to the Tweet Repository, exposing Tweet statistics HTTP endpoints
* **Frontend:** A simple dashboard frontend displaying info from Stats API