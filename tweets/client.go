package tweets

import(
	"net/http"
)

// We wrap the http client to ease unit testing
type HttpClient interface{
	Do(req *http.Request) (*http.Response, error)
	Url() string
}

type MockClient struct{
	DoFunc func(req *http.Request) (*http.Response, error)
	UrlFunc func() string
}

type TweetClient struct{
	Client http.Client
}

// The normal tweet client delegates 'Do' to the underlying http.Client
func (tc TweetClient) Do(req *http.Request) (*http.Response, error) {
	return tc.Client.Do(req)
}

// Url is abstracted to ease unit testing
func (tc TweetClient) Url() string {
	return "https://api.twitter.com/2/tweets/sample/stream?tweet.fields=text,id,created_at,author_id,entities,lang,public_metrics"
}

func (m MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func (m MockClient) Url() string {
	return m.UrlFunc()
}

func NewTweetClient() TweetClient{
	return TweetClient{
		Client: http.Client{},
	}
}
