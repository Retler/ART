package tweets

// Some mock tweets for testing purposes
var MockTweet1 = Tweet{
	Data: Data{
		TweetID:   "1396383361833209856",
		Content:   "RT @momy9775: What an awesome tweet this is!",
		AuthorID:  "1085741751174721536",
		CreatedAt: "2021-05-23T08:31:51.000Z",
		Language:  "en",
		PublicMetrics: PublicMetrics{
			RetweetCount: 12927,
			LikeCount:    0,
		},
		Entities: Entities{
			Hashtags: []Hashtag{
				Hashtag{
					Tag: "Test",
				},
			},
		},
	},
	Sentiment: 0.5,
}

var MockTweet2 = Tweet{
	Data: Data{
		TweetID:   "1396383361833209857",
		Content:   "This is alsow a nice tweet",
		AuthorID:  "1085741751174721536",
		CreatedAt: "2021-04-23T09:31:51.000Z",
		Language:  "en",
		PublicMetrics: PublicMetrics{
			RetweetCount: 12927,
			LikeCount:    0,
		},
		Entities: Entities{
			Hashtags: []Hashtag{
				Hashtag{
					Tag: "BLM",
				},
			},
		},
	},
	Sentiment: 0.5,
}

var MockTweet3 = Tweet{
	Data: Data{
		TweetID:   "1396383361833209858",
		Content:   "Man I hate ice cream.. #IceCreamHatersUnite",
		AuthorID:  "1085741751174721536",
		CreatedAt: "2021-03-23T08:31:51.000Z",
		Language:  "en",
		PublicMetrics: PublicMetrics{
			RetweetCount: 12927,
			LikeCount:    0,
		},
		Entities: Entities{
			Hashtags: []Hashtag{
				Hashtag{
					Tag: "IceCreamHatersUnite",
				},
			},
		},
	},
	Sentiment: -0.65,
}

var MockTweet4 = Tweet{
	Data: Data{
		TweetID:   "1496383361833209810",
		Content:   "I also hate ice cream!! #IceCreamHatersUnite",
		AuthorID:  "1085741751174721599",
		CreatedAt: "2021-03-23T08:31:25.000Z",
		Language:  "en",
		PublicMetrics: PublicMetrics{
			RetweetCount: 100,
			LikeCount:    0,
		},
		Entities: Entities{
			Hashtags: []Hashtag{
				Hashtag{
					Tag: "IceCreamHatersUnite",
				},
			},
		},
	},
	Sentiment: -0.45,
}
