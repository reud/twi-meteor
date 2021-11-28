package adapter

import (
	"github.com/g8rswimmer/go-twitter"
	"github.com/michimani/gotwi/resources"
	"github.com/reud/twi-meteor/domain"
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
func V1ToDomainTweetModel(obj resources.Tweet) (*domain.Tweet, error) {
	return &domain.Tweet{
		CreatedAt: *obj.CreatedAt,
		ID:        *obj.ID,
		Text:      *obj.Text,
	}, nil
}

func ToDomainLikeDataModel(dirty resources.User) domain.LikeData {
	return domain.LikeData{
		ID:       *dirty.ID,
		Name:     *dirty.Name,
		Username: *dirty.Username,
	}
}
