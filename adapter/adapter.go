package adapter

import (
	"github.com/michimani/gotwi/resources"
	"github.com/reud/twi-meteor/domain"
)

// ToDomainTweetModel は infra(v1) -> domainへのコンバータ
func ToDomainTweetModel(obj resources.Tweet) (*domain.Tweet, error) {
	var rts []domain.ReferencedTweet
	for _, tweet := range obj.ReferencedTweets {
		r, e := ToDomainReferencedTweetModel(tweet)
		if e != nil {
			return nil, e
		}
		rts = append(rts, r)
	}

	return &domain.Tweet{
		CreatedAt:       *obj.CreatedAt,
		ID:              *obj.ID,
		Text:            *obj.Text,
		ReferencedTweet: rts,
	}, nil
}

func ToDomainReferencedTweetModel(obj resources.ReferencedTweet) (domain.ReferencedTweet, error) {
	return domain.ReferencedTweet{
		Type: *obj.Type,
		ID:   *obj.ID,
	}, nil
}

func ToDomainLikeDataModel(dirty resources.User) domain.LikeData {
	return domain.LikeData{
		ID:       *dirty.ID,
		Name:     *dirty.Name,
		Username: *dirty.Username,
	}
}
