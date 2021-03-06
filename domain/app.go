package domain

import (
	"strings"
	"time"
)

// ClockInterface time.Now()のモックにするためのインターフェイス
type ClockInterface interface {
	Now() time.Time
}

type AdapterInterface interface {
	FetchTweets() ([]Tweet, error)
	LikingUsers(tweetID string) ([]LikeData, error)
	DestroyTweet(tweetID string) error
}

type ApplicationInterface interface {
	CheckDeletableTweet(tweet Tweet) (isOK bool, err error)
	DestroyTweet(tweet Tweet) (err error)
}

type Application struct {
	Client      AdapterInterface
	MyTwitterID string
	clock       ClockInterface
}

type Clock struct{}

func (c *Clock) Now() time.Time {
	return time.Now()
}

func GenApplication(client AdapterInterface, myTwitterID string, clock ClockInterface) ApplicationInterface {
	return &Application{
		Client:      client,
		MyTwitterID: myTwitterID,
		clock:       clock,
	}
}

// CheckDeletableTweet はそのツイートが削除してよいか判定する。
// isOK 判定が上手くいき、そのツイートが削除して良い時はtrue それ以外false
// err 判定が上手くいかない時や、削除NGの時返す。
func (app *Application) CheckDeletableTweet(tweet Tweet) (isOK bool, err error) {
	// 現在時刻より24時間以上前のツイートでない場合は終了
	if !tweet.CreatedAt.Before(app.clock.Now().Add(-time.Hour * 24)) {
		return false, &CheckFailedError{
			Message: "24時間経っていないツイートです",
		}
	}

	// これは誰かへのリプライのためスキップ(暫定的)
	if strings.Contains(tweet.Text, "@") {
		return false, &CheckFailedError{Message: "誰かへのリプライです。"}
	}

	likeData, err := app.Client.LikingUsers(tweet.ID)
	if err != nil {
		return false, err
	}

	for _, likeUser := range likeData {
		if likeUser.ID == app.MyTwitterID {
			return false, &CheckFailedError{Message: "あなたがいいねしたツイートです"}
		}
	}

	// 誰かからリプライがきている時にスキップ
	for _, t := range tweet.ReferencedTweet {
		if t.Type == "replied_to" {
			return false, &CheckFailedError{Message: "リプライが来ています。"}
		}
	}

	return true, nil
}

// DestroyTweet はそのツイートを削除する。
func (app *Application) DestroyTweet(tweet Tweet) (err error) {
	return app.Client.DestroyTweet(tweet.ID)
}
