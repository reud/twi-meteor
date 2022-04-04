package infra

import (
	"context"
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/resources"
	"github.com/michimani/gotwi/tweets"
	"github.com/michimani/gotwi/tweets/types"
	"os"
	"strings"
)

const LOOPS = 10 // ツイート取得APIを一度の探索で叩く回数

/*
	verify_credentialsを使うと依存ライブラリが増えるため、.envに直接usernameを入れる方向で一旦対応
*/
type TwitterClientInterface interface {
	DestroyTweet(tweetID string) error
	LikingUsers(tweetID string) ([]resources.User, error)
	FetchMyTweets() ([]resources.Tweet, error)
	FetchMyTweetsOnce(paginationToken string) ([]resources.Tweet, string, error)
}

type TwitterClient struct {
	client    *gotwi.GotwiClient
	twitterID string
}

func GenTwitterClient(at, ats, ck, cs, twitterID string) (TwitterClientInterface, error) {
	os.Setenv("GOTWI_API_KEY", ck)
	os.Setenv("GOTWI_API_KEY_SECRET", cs)
	in := &gotwi.NewGotwiClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           at,
		OAuthTokenSecret:     ats,
	}
	client, err := gotwi.NewGotwiClient(in)
	if err != nil {
		return nil, err
	}
	return TwitterClient{client: client, twitterID: twitterID}, nil
}

func (tc TwitterClient) LikingUsers(tweetID string) ([]resources.User, error) {
	p := &types.TweetLikesLikingUsersParams{
		ID: tweetID,
	}
	res, err := tweets.TweetLikesLikingUsers(context.Background(), tc.client, p)
	if err != nil {
		return nil, err
	}
	if res.HasPartialError() {
		return nil, combineError(res.Errors)
	}

	return res.Data, nil
}

func (tc TwitterClient) DestroyTweet(tweetID string) error {
	p := &types.ManageTweetsDeleteParams{ID: tweetID}
	res, err := tweets.ManageTweetsDelete(context.Background(), tc.client, p)
	if err != nil {
		return err
	}
	if !*res.Data.Deleted {
		return fmt.Errorf("failed to delete tweets(but twitter has no say)")
	}
	return nil
}

// FetchMyTweetsOnce は一回だけ取る
func (tc TwitterClient) FetchMyTweetsOnce(paginationToken string) ([]resources.Tweet, string, error) {
	p := &types.TweetTimelinesTweetsParams{
		ID:              tc.twitterID,
		PaginationToken: paginationToken,
		TweetFields:     fields.TweetFieldList{fields.TweetFieldCreatedAt, fields.TweetFieldReferencedTweets},
		MaxResults:      100,
	}
	res, err := tweets.TweetTimelinesTweets(context.Background(), tc.client, p)
	if err != nil {
		return nil, "", err
	}
	if len(res.Errors) != 0 {
		return nil, "", combineError(res.Errors)
	}
	return res.Data, p.PaginationToken, err
}

// FetchMyTweets は与えられた回数分取る
func (tc TwitterClient) FetchMyTweets() ([]resources.Tweet, error) {
	res, paginationToken, err := tc.FetchMyTweetsOnce("")
	if err != nil {
		return nil, err
	}
	result := res
	for i := 0; i < LOOPS; i++ {
		if paginationToken == "" {
			return result, nil
		}
		res, newPaginationToken, err := tc.FetchMyTweetsOnce(paginationToken)
		if err != nil {
			return nil, err
		}
		result = append(result, res...)
		paginationToken = newPaginationToken
	}
	return result, nil
}

func combineError(errors []resources.PartialError) error {
	var errstrings []string
	for _, e := range errors {
		errstrings = append(errstrings, pp.Sprint(e))
	}
	return fmt.Errorf(strings.Join(errstrings, "\n"))
}
