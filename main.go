package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/reud/twi-meteor/adapter"
	"github.com/reud/twi-meteor/domain"
	"github.com/reud/twi-meteor/env"
	"github.com/reud/twi-meteor/infra"
	"log"
)

func main() {
	err := godotenv.Overload()
	if err != nil {
		log.Printf("failed to load .env")
	}

	con := env.GetTwitterConfig()

	infraCl, err := infra.GenTwitterClient(con.AccessToken, con.AccessTokenSecret, con.ConsumerKey, con.ConsumerSecret, con.TwitterID)
	if err != nil {
		log.Fatal(err)
		return
	}

	adapterCl := adapter.GenAdapterClient(infraCl)
	tweets, err := adapterCl.FetchTweets()

	if err != nil {
		log.Fatal(err)
	}

	app := domain.GenApplication(adapterCl, con.TwitterID, &domain.Clock{})
	for _, tweet := range tweets {
		isOK, err := app.CheckDeletableTweet(tweet)
		if err != nil {
			switch err := err.(type) {
			case *domain.CheckFailedError:
				_, err2 := pp.Print(*err)
				if err2 != nil {
					log.Fatal(err2)
				}
				continue
			default:
				log.Fatal(err)
			}
		}

		fmt.Printf("tweet (deleted: %+v): \n", isOK)
		fmt.Print("deleting...")
		if err = app.DestroyTweet(tweet); err != nil {
			log.Fatal(err)
		}
		if _, err = pp.Print(tweet); err != nil {
			log.Fatal(err)
		}
	}
}
