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
	"strings"
)

/*
	verify_credentialsを使うと依存ライブラリが増えるため、.envに直接usernameを入れる方向で一旦対応
*/
type TwitterClientInterface interface {
	LikingUsers(tweetID string) ([]resources.User, error)
	FetchMyTweets(paginationToken string) ([]resources.Tweet, error)
}

type TwitterClient struct {
	client    *gotwi.GotwiClient
	twitterID string
}

func GenTwitterClient(at string, ats string, twitterID string) (TwitterClientInterface, error) {
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

func (tc TwitterClient) FetchMyTweets(paginationToken string) ([]resources.Tweet, string, error) {
	p := &types.TweetTimelinesTweetsParams{
		ID:              tc.twitterID,
		PaginationToken: paginationToken,
		TweetFields:     fields.TweetFieldList{fields.TweetFieldCreatedAt},
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

func combineError(errors []resources.PartialError) error {
	var errstrings []string
	for _, e := range errors {
		errstrings = append(errstrings, pp.Sprint(e))
	}
	return fmt.Errorf(strings.Join(errstrings, "\n"))
}
