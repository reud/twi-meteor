package infra

import (
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/reud/twi-meteor/env"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Setup(t *testing.T) env.TwitterConfig {
	err := godotenv.Overload("../.env")
	if err != nil {
		log.Printf("failed to load .env")
	}
	evv := env.GetTwitterConfig()
	t.Setenv("GOTWI_API_KEY", evv.ConsumerKey)
	t.Setenv("GOTWI_API_KEY_SECRET", evv.ConsumerSecret)
	return evv
}

func TestTwitterClient_FetchMyTweetsOnce(t *testing.T) {
	config := Setup(t)
	client, err := GenTwitterClient(config.AccessToken, config.AccessTokenSecret, config.ConsumerKey, config.ConsumerSecret, config.TwitterID)
	if err != nil {
		assert.Nil(t, err)
	}

	tweets, _, err := client.FetchMyTweetsOnce("")
	assert.Nil(t, err)
	_, err = pp.Print(tweets)
	assert.Nil(t, err)
}

func TestTwitterClient_LikingUsers(t *testing.T) {
	config := Setup(t)
	client, err := GenTwitterClient(config.AccessToken, config.AccessTokenSecret, config.ConsumerKey, config.ConsumerSecret, config.TwitterID)

	if err != nil {
		assert.Nil(t, err)
	}

	users, err := client.LikingUsers("1464805692355735557")
	assert.Nil(t, err)
	_, err = pp.Print(users)
	assert.Nil(t, err)
}

func TestTwitterClient_FetchMyTweets(t *testing.T) {
	config := Setup(t)
	client, err := GenTwitterClient(config.AccessToken, config.AccessTokenSecret, config.ConsumerKey, config.ConsumerSecret, config.TwitterID)

	assert.Nil(t, err)
	tweets, err := client.FetchMyTweets()
	assert.Nil(t, err)
	_, err = pp.Print(tweets)
	assert.Nil(t, err)
}
