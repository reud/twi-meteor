package v1

import (
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/reud/twi-meteor/env"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

/*
	infra層のテストは探索型テスト
	実際にAPIにアクセスしてコードの仕様を確認する
*/

func Setup() TwitterV1ClientInterface {
	err := godotenv.Overload("../../.env")
	if err != nil {
		log.Printf("failed to load .env")
	}
	con := env.GetTwitterConfig()
	v1client := GenTwitterV1Client(V1Config{
		ConsumerKey:       con.ConsumerKey,
		ConsumerSecret:    con.ConsumerSecret,
		AccessToken:       con.AccessToken,
		AccessTokenSecret: con.AccessTokenSecret,
	})
	return v1client
}

func TestTwitterV1Client_FetchMyTweets(t *testing.T) {
	client := Setup()
	tweets, err := client.FetchMyTweets(nil)
	assert.Nil(t, err)
	pp.Print(tweets)
}
