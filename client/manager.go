package client

import (
	"fmt"
	"github.com/g8rswimmer/go-twitter"
	v1 "github.com/reud/twi-meteor/infra/v1"
	v2 "github.com/reud/twi-meteor/infra/v2"
)

const LOOPS = 10

type Client struct {
	v1client *v1.TwitterV1Client
	v2client *v2.TwitterV2Client
}

func GenClient(v1 *v1.TwitterV1Client, v2 *v2.TwitterV2Client) *Client {
	return &Client{
		v1client: v1,
		v2client: v2,
	}
}

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
