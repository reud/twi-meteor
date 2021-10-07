package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/g8rswimmer/go-twitter"
	"net/http"
)

type V2Config struct {
	BearerToken string
	TwitterID   string
}

type TwitterV2Client struct {
	BearerToken string
	TwitterID   string
}

func GenTwitterV2Client(con V2Config) *TwitterV2Client {
	return &TwitterV2Client{
		BearerToken: con.BearerToken,
		TwitterID:   con.TwitterID,
	}
}

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func (v2 TwitterV2Client) FetchMyTweets() {
	user := &twitter.User{
		Authorizer: authorize{
			Token: v2.BearerToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	tweetOpts := twitter.UserTimelineOpts{
		TweetFields: []twitter.TweetField{
			twitter.TweetFieldAuthorID,
			twitter.TweetFieldContextAnnotations,
			twitter.TweetFieldCreatedAt,
			twitter.TweetFieldID,
			twitter.TweetFieldPublicMetrics,
			twitter.TweetFieldText,
		},
		UserFields: []twitter.UserField{
			twitter.UserFieldCreatedAt,
			twitter.UserFieldDescription,
			twitter.UserFieldName,
			twitter.UserFieldURL,
			twitter.UserFieldUserName,
		},
		Expansions:  []twitter.Expansion{},
		PlaceFields: []twitter.PlaceField{},
		PollFields:  []twitter.PollField{},
		MaxResults:  10,
	}

	userTweets, err := user.Tweets(context.Background(), v2.TwitterID, tweetOpts)
	var tweetErr *twitter.TweetErrorResponse
	switch {
	case errors.As(err, &tweetErr):
		printTweetError(tweetErr)
	case err != nil:
		fmt.Println(err)
	default:
		printUserTweets(userTweets)
	}
}

func printUserTweets(userTweets *twitter.UserTimeline) {
	enc, err := json.MarshalIndent(userTweets, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enc))
}

func printTweetLookup(lookup twitter.TweetLookup) {
	enc, err := json.MarshalIndent(lookup, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enc))
}

func printTweetError(tweetErr *twitter.TweetErrorResponse) {
	enc, err := json.MarshalIndent(tweetErr, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enc))
}
