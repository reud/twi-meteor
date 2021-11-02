package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/reud/twi-meteor/adapter"
	"github.com/reud/twi-meteor/domain"
	"github.com/reud/twi-meteor/env"
	"github.com/reud/twi-meteor/infra"
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
	strTwitterId := strconv.FormatInt(id, 10)
	if err != nil {
		log.Fatal(err)
	}

	v2client := v2.GenTwitterV2Client(v2.V2Config{
		BearerToken: con.BearerToken,
		TwitterID:   strTwitterId,
	})

	infraCl := infra.GenClient(v1cleint, v2client)
	adapterCl := adapter.GenAdapterClient(infraCl)
	tweets, err := adapterCl.FetchTweets()

	if err != nil {
		log.Fatal(err)
	}

	app := domain.GenApplication(adapterCl, strTwitterId)
	for _, tweet := range tweets {
		isOK, err := app.CheckDeletableTweet(tweet)
		if err != nil {
			switch err := err.(type) {
			case *domain.CheckFailedError:
				_, err2 := pp.Print(*err)
				if err2 != nil {
					log.Fatal(err2)
				}
			default:
				log.Fatal(err)
			}

		}
		fmt.Printf("tweet (deleted: %+v): \n", isOK)
		_, err = pp.Print(tweet)
		if err != nil {
			log.Fatal(err)
		}
	}
}
