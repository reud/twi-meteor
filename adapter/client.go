package adapter

import (
	"github.com/reud/twi-meteor/domain"
	"github.com/reud/twi-meteor/infra"
	"strconv"
	"time"
)

// GenAdapterClient infra -> domainへのデータを成形するクライアントを生成する。
func GenAdapterClient(infraClient infra.ClientInterface) *Client {
	return &Client{infraClient}
}

type Client struct {
	InfraClient infra.ClientInterface
}

// FetchTweets はinfra層の同メソッドを呼び出し、domain層のモデルに変換して返す
func (c *Client) FetchTweets() ([]domain.Tweet, error) {
	time.Sleep(time.Second) // for 429 対策

	infraResult, err := c.InfraClient.FetchTweets()
	if err != nil {
		return nil, err
	}

	var domainResult []domain.Tweet

	for _, infraTweet := range infraResult {
		domainTweet, err := ToDomainTweetModel(infraTweet)
		if err != nil {
			return nil, err
		}

		domainResult = append(domainResult, *domainTweet)
	}

	return domainResult, nil
}

// LikingUsers はinfra層の同メソッドを呼び出し、domain層のモデルに変換して返す
func (c *Client) LikingUsers(tweetID string) ([]domain.LikeData, error) {
	time.Sleep(time.Second) // for 429 対策
	infraLikeData, err := c.InfraClient.LikingUsers(tweetID)
	if err != nil {
		return nil, err
	}

	var domainLikeData []domain.LikeData

	for _, infraLikeDataElement := range infraLikeData {
		domainLikeDataElement := ToDomainLikeDataModel(infraLikeDataElement)
		domainLikeData = append(domainLikeData, domainLikeDataElement)
	}
	return domainLikeData, nil
}

// DestroyTweet はinfra層の同メソッドを呼び出す
func (c *Client) DestroyTweet(tweetID string) error {
	time.Sleep(time.Second) // for 429 対策
	tweetID64, err := strconv.ParseInt(tweetID, 10, 64)
	if err != nil {
		return err
	}
	_, err = c.InfraClient.DestroyTweet(tweetID64)
	return err
}
