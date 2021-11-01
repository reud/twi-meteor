package v1

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/k0kubun/pp"
	"log"
)

type V1Config struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type TwitterV1ClientInterface interface {
	LookupID() (int64, error)
	DestroyTweet(tweetID int64) (int64, error)
}

type TwitterV1Client struct {
	client *twitter.Client
}

func GenTwitterV1Client(con V1Config) TwitterV1ClientInterface {
	config := oauth1.NewConfig(con.ConsumerKey, con.ConsumerSecret)
	token := oauth1.NewToken(con.AccessToken, con.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return &TwitterV1Client{client: client}
}

func (v1 *TwitterV1Client) LookupID() (int64, error) {
	user, _, err := v1.client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

func (v1 *TwitterV1Client) DestroyTweet(tweetID int64) (int64, error) {
	deletedTweet, _, err := v1.client.Statuses.Destroy(tweetID, &twitter.StatusDestroyParams{})
	if err != nil {
		log.Printf("DELETE FAILED】\n err: %+v", pp.Sprint(err))
		return 0, err
	}
	log.Printf("【DELETED】 \n %+v", pp.Sprint(deletedTweet))
	return deletedTweet.ID, err
}
