package main

import (
	"github.com/joho/godotenv"
	"github.com/reud/twi-meteor/env"
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

	v2client.FetchMyTweets()
}
