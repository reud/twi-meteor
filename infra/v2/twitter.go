package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/g8rswimmer/go-twitter"
	"io/ioutil"
	"net/http"
)

type LikeData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type LikingUsersResponse struct {
	Data []LikeData `json:"data"`
}

type V2Config struct {
	BearerToken string
	TwitterID   string
}

type TwitterV2Client struct {
	BearerToken string
	TwitterID   string
}

func GenTwitterV2Client(con V2Config) TwitterV2ClientInterface {
	return &TwitterV2Client{
		BearerToken: con.BearerToken,
		TwitterID:   con.TwitterID,
	}
}

type TwitterV2ClientInterface interface {
	TweetLookup(tweetID string) (*twitter.TweetLookup, error)
	FetchMyTweets(paginationToken string) ([]twitter.TweetObj, string, error)
	LikingUsers(tweetID string) ([]LikeData, error)
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

// LikingUsers はtweetIDにいいねしたユーザの情報を取得する。ライブラリが無いのでスクラッチで書いている
func (v2 TwitterV2Client) LikingUsers(tweetID string) ([]LikeData, error) {
	url := fmt.Sprintf("https://api.twitter.com/2/tweets/%s/liking_users", tweetID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []LikeData{}, err
	}
	req.Header.Set("Authorization", "Bearer "+v2.BearerToken)
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return []LikeData{}, err
	}

	if resp.StatusCode != 200 {
		return []LikeData{}, fmt.Errorf("failed to request, code: %+v \n status: %+v", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []LikeData{}, err
	}
	var decoded LikingUsersResponse
	if err := json.Unmarshal(byteArray, &decoded); err != nil {
		return []LikeData{}, err
	}
	return decoded.Data, nil
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
