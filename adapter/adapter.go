package adapter

import (
	"github.com/g8rswimmer/go-twitter"
	"github.com/reud/twi-meteor/domain"
	v2 "github.com/reud/twi-meteor/infra/v2"
	"time"
)

// ToDomainTweetModel は infra -> domainへのコンバータ
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
func ToDomainLikeDataModel(dirty v2.LikeData) domain.LikeData {
	return domain.LikeData{
		ID:       dirty.ID,
		Name:     dirty.Name,
		Username: dirty.Username,
	}
}
