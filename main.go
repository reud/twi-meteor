package main

import (
	"github.com/joho/godotenv"
	"github.com/reud/twi-meteor/client"
	"github.com/reud/twi-meteor/env"
	"github.com/reud/twi-meteor/infra/http"
	v1 "github.com/reud/twi-meteor/infra/v1"
	v2 "github.com/reud/twi-meteor/infra/v2"
	"log"
	"strconv"
)

func main() {
	err := godotenv.Overload()
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
	if err != nil {
		log.Fatal(err)
	}

	v2client := v2.GenTwitterV2Client(v2.V2Config{
		BearerToken: con.BearerToken,
		TwitterID:   strconv.FormatInt(id, 10),
	})

	cl := client.GenClient(v1cleint, v2client)
	twts, err := cl.FetchTweets()
	if err != nil {
		log.Fatal(err)
	}
	for _, tweet := range twts {
		likes := http.LikingUsers(tweet.ID, con.BearerToken)
		found := false
		for _, user := range likes {
			if user.ID == strconv.FormatInt(id, 10) {
				log.Printf("saved: %+v", tweet.ID)
				found = true
				break
			}
		}
		if !found {
			convertedStrInt64, err := strconv.ParseInt(tweet.ID, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			_, err = v1cleint.DestroyTweet(convertedStrInt64)
			if err != nil {
				log.Print(err)
			}
		}
	}
}
