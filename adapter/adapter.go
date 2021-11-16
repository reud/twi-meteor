package adapter

import (
	v1client "github.com/dghubble/go-twitter/twitter"
	"github.com/g8rswimmer/go-twitter"
	"github.com/reud/twi-meteor/domain"
	v2 "github.com/reud/twi-meteor/infra/v2"
	"time"
)

// ToDomainTweetModel は infra(v2) -> domainへのコンバータ
func ToDomainTweetModel(obj twitter.TweetObj) (*domain.Tweet, error) {
	tweetTime, err := time.Parse("2006-01-02T15:04:05.000Z", obj.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Tweet{
		CreatedAt: tweetTime,
		ID:        obj.ID,
		Text:      obj.Text,
	}, nil
}

// V1ToDomainTweetModel は infra(v1) -> domainへのコンバータ
func V1ToDomainTweetModel(obj v1client.Tweet) (*domain.Tweet, error) {
	tweetTime, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", obj.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Tweet{
		CreatedAt: tweetTime,
		ID:        obj.IDStr,
		Text:      obj.Text,
	}, nil
}

func ToDomainLikeDataModel(dirty v2.LikeData) domain.LikeData {
	return domain.LikeData{
		ID:       dirty.ID,
		Name:     dirty.Name,
		Username: dirty.Username,
	}
}
