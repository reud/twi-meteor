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

func (v2 TwitterV2Client) TweetLookup(tweetID string) (*twitter.TweetLookup, error) {
	tweet := &twitter.Tweet{
		Authorizer: authorize{
			Token: v2.BearerToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	fieldOpts := twitter.TweetFieldOptions{
		Expansions:  []twitter.Expansion{twitter.ExpansionEntitiesMentionsUserName, twitter.ExpansionAuthorID},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldAttachments},
	}

	lookups, err := tweet.Lookup(context.Background(), []string{tweetID}, fieldOpts)
	var tweetErr *twitter.TweetErrorResponse
	switch {
	case errors.As(err, &tweetErr):
		printTweetError(tweetErr)
		return nil, tweetErr
	case err != nil:
		fmt.Println(err)
		return nil, err
	default:
		x := lookups[tweetID]
		return &x, nil
	}
}

func (v2 TwitterV2Client) FetchMyTweets(paginationToken string) ([]twitter.TweetObj, string, error) {
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
		Expansions:      []twitter.Expansion{},
		PlaceFields:     []twitter.PlaceField{},
		PollFields:      []twitter.PollField{},
		MaxResults:      100,
		PaginationToken: paginationToken,
	}

	userTweets, err := user.Tweets(context.Background(), v2.TwitterID, tweetOpts)
	var tweetErr *twitter.TweetErrorResponse
	switch {
	case errors.As(err, &tweetErr):
		printTweetError(tweetErr)
		return []twitter.TweetObj{}, "", tweetErr
	case err != nil:
		fmt.Println(err)
		return []twitter.TweetObj{}, "", err
	default:
		printUserTweets(userTweets)
		return userTweets.Tweets, userTweets.Meta.NextToken, nil
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
