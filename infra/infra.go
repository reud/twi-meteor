package infra

import (
	"fmt"
	"github.com/g8rswimmer/go-twitter"
	v1 "github.com/reud/twi-meteor/infra/v1"
	v2 "github.com/reud/twi-meteor/infra/v2"
)

type ClientInterface interface {
	FetchTweets() ([]twitter.TweetObj, error)
	LikingUsers(tweetID string) ([]v2.LikeData, error)
	DestroyTweet(tweetID int64) (int64, error)
}

const LOOPS = 10

type Client struct {
	v1client v1.TwitterV1ClientInterface
	v2client v2.TwitterV2ClientInterface
}

func GenClient(v1 v1.TwitterV1ClientInterface, v2 v2.TwitterV2ClientInterface) ClientInterface {
	return &Client{
		v1client: v1,
		v2client: v2,
	}
}

// FetchTweets はtwitterから自身のツイートを取得する。
func (c *Client) FetchTweets() ([]twitter.TweetObj, error) {
	tweet, next, err := c.v2client.FetchMyTweets("")
	if err != nil {
		fmt.Printf("%+v", err)
	}
	result := tweet
	for i := 0; i < LOOPS; i++ {
		tweet, next, err = c.v2client.FetchMyTweets(next)
		result = append(result, tweet...)
		if err != nil {
			fmt.Printf("%+v th error occured: %+v", i, err)
		}
		if len(next) == 0 {
			break
		}
	}
	return result, nil
}

// LikingUsers はツイートからからいいねしたユーザを取得する。
func (c *Client) LikingUsers(tweetID string) ([]v2.LikeData, error) {
	return c.v2client.LikingUsers(tweetID)
}

// DestroyTweet はツイートの削除を行う
func (c *Client) DestroyTweet(tweetID int64) (int64, error) {
	return c.v1client.DestroyTweet(tweetID)
}
