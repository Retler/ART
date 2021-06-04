package tweets

import (
	"fmt"
	"testing"
	"time"
)

// Test that all hashtags are aggregated correctly
func TestTopHashtags(t *testing.T) {
	tweets := Tweets{
		[]Tweet{
			MockTweet1,
			MockTweet2,
			MockTweet3,
			MockTweet4,
		},
	}

	hashtags := tweets.TopXHashtags(100)

	if len(hashtags) != 3 {
		t.Errorf("Mock tweets should only contain 3 hashtags but got: %d", len(hashtags))
	}

	// Returned hashtags should be sorted by their count.
	if hashtags[0].Hashtag != "IceCreamHatersUnite" {
		t.Errorf("Wrong hashtag at 1st pos in list. Expected: %s but got %s", "IceCreamHatersUnite", hashtags[0].Hashtag)
	}
}

// Test that sentiment is computed correctly
func TestSentiment(t *testing.T) {
	tweets := Tweets{
		[]Tweet{
			MockTweet1,
			MockTweet2,
			MockTweet3,
			MockTweet4,
		},
	}

	sentiment := fmt.Sprintf("%.2f", tweets.Sentiment())
	expectedSentiment := "-0.03"

	// Returned hashtags should be sorted by their count.
	if sentiment != expectedSentiment {
		t.Errorf("Expected sentiment %s but got %s", expectedSentiment, sentiment)
	}
}

// Test that grouping by hour and minute works correctly
func TestGroupTweets(t *testing.T) {
	tweets := Tweets{
		[]Tweet{
			MockTweet1,
			MockTweet2,
			MockTweet3,
			MockTweet4,
		},
	}

	groupsMinute, err := tweets.ByDuration(time.Minute)
	if err != nil {
		t.Errorf("Error when grouping by minute: %s", err)
	}

	groupsHour, err := tweets.ByDuration(time.Hour)
	if err != nil {
		t.Errorf("Error when grouping by hour: %s", err)
	}

	expectedGroups := 3
	if len(groupsMinute) != expectedGroups {
		t.Errorf("Expected minute groups: %d but got: %d", expectedGroups, len(groupsMinute))
	}
	if len(groupsHour) != expectedGroups {
		t.Errorf("Expected hour groups: %d but got: %d", expectedGroups, len(groupsHour))
	}
}
