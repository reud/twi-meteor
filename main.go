package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

type TwitterConfig struct {
	ConsumerKey string `envconfig:"CONSUMER_KEY" required:"true"`
	ConsumerSecret string `envconfig:"CONSUMER_SECRET" required:"true"`
	AccessToken string `envconfig:"ACCESS_TOKEN" required:"true"`
	AccessTokenSecret string `envconfig:"ACCESS_TOKEN_SECRET" required:"true"`
}

func getTwitterConfig() TwitterConfig {
	var con TwitterConfig
	err := envconfig.Process("",&con)
	if err != nil {
		panic(err)
	}
	return con
}

func genTwitterClient(con TwitterConfig) *twitter.Client {
	config := oauth1.NewConfig(con.ConsumerKey, con.ConsumerSecret)
	token := oauth1.NewToken(con.AccessToken, con.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return client
}

// 自分が自分のツイートをいいねしている時setに入る。
// 雑に自分のいいね履歴から2000個ぐらい持ってくる。 twitter v2ならもっと楽だった・・・
func genSelfFavoriteMap(tw *twitter.Client,userID int64) map[int64]struct{} {
	m := map[int64]struct{}{}

	// 呼び出し回数
	calls := 10
	sinceId := int64(0)

	for i := 0; i < calls; i++  {
		time.Sleep(time.Second * 10)
		log.Printf("fetch favorite %+v\n",i)
		tweets, _ , err := tw.Favorites.List(&twitter.FavoriteListParams{
			UserID:          userID,
			Count:           200,
			SinceID:         sinceId,
			MaxID:           0,
			IncludeEntities: nil,
			TweetMode:       "",
		})
		if err != nil {
			panic(err)
		}
		if len(tweets) == 0 {
			log.Printf("no favorite tweets\n")
			break
		}
		sinceId = tweets[len(tweets)-1].ID

		for _, tweet := range tweets {
			if tweet.User.ID == userID {
				m[tweet.ID] = struct{}{}
			}
		}
	}
	return m
}

func main() {
	err := godotenv.Overload()
	if err != nil {
		log.Printf("failed to load .env")
	}

	con := getTwitterConfig()
	client := genTwitterClient(con)
	user,_,err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		log.Fatalf("cant fetch user error %+v %+v",pp.Sprint(err),pp.Sprint(con))
	}

	r,httpRes,err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID:          user.ID,
		Count: 200,
	})

	if err != nil {
		log.Fatalf("failed to fetch tweet: %+v %+v",httpRes,err)
		return
	}

	favoritesMyTweetMap := genSelfFavoriteMap(client,user.ID)
	for idx, tweet  := range r {
		if _, ok := favoritesMyTweetMap[tweet.ID]; ok {
			log.Printf("idx: %+v 【IGNORED】 \n %+v",idx,pp.Sprint(tweet))
		} else {
			deletedTweet,_,err := client.Statuses.Destroy(tweet.ID,&twitter.StatusDestroyParams{})
			if err != nil {
				log.Printf("idx: %+v 【DELETE FAILED】\n err: %+v \n tweet: %+v",pp.Sprint(err),pp.Sprint(tweet))
			}
			log.Printf("idx: %+v 【DELETED】 \n %+v",idx,pp.Sprint(deletedTweet))
		}
	}
}