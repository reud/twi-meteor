package v2

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/reud/twi-meteor/env"
	v1 "github.com/reud/twi-meteor/infra/v1"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
)

/*
	infra層のテストは探索型テスト
	実際にAPIにアクセスしてコードの仕様を確認する
*/

func Setup() TwitterV2ClientInterface {
	err := godotenv.Overload("../../.env")
	if err != nil {
		log.Printf("failed to load .env")
	}
	con := env.GetTwitterConfig()
	v1cleint := v1.GenTwitterV1Client(v1.V1Config{
		ConsumerKey:       con.ConsumerKey,
		ConsumerSecret:    con.ConsumerSecret,
		AccessToken:       con.AccessToken,
		AccessTokenSecret: con.AccessTokenSecret,
	})

	id, err := v1cleint.LookupID()
	strTwitterId := strconv.FormatInt(id, 10)
	if err != nil {
		log.Fatal(err)
	}
	return GenTwitterV2Client(V2Config{
		BearerToken: con.BearerToken,
		TwitterID:   strTwitterId,
	})
}

func TestTwitterV2Client_FetchMyTweets(t *testing.T) {
	v2Client := Setup()
	t.Run("データが適切に取得できること", func(t *testing.T) {
		result, next, err := v2Client.FetchMyTweets("")
		fmt.Printf("next: %+v", next)
		assert.Nil(t, err)
		for _, tweet := range result {
			_, err := pp.Print(tweet)
			assert.Nil(t, err)
		}
	})
}

func TestTwitterV2Client_LikingUsers(t *testing.T) {
	v2Client := Setup()
	t.Run("データが適切に取得できること", func(t *testing.T) {
		result, err := v2Client.LikingUsers("1455116170864979975")
		assert.Nil(t, err)

		for _, tweet := range result {
			_, err := pp.Print(tweet)
			assert.Nil(t, err)
		}
	})
}
