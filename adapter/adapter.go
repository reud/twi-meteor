package adapter

import (
	"github.com/michimani/gotwi/resources"
	"github.com/reud/twi-meteor/domain"
)

// ToDomainTweetModel は infra(v1) -> domainへのコンバータ
func ToDomainTweetModel(obj resources.Tweet) (*domain.Tweet, error) {
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
