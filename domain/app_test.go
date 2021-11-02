package domain

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type AdapterMock struct {
	LikingUsersResult []LikeData
}

func (ad *AdapterMock) FetchTweets() ([]Tweet, error) {
	return nil, nil
}

func (ad *AdapterMock) LikingUsers(tweetID string) ([]LikeData, error) {
	return ad.LikingUsersResult, nil
}

func (ad *AdapterMock) DestroyTweet(tweetID string) error {
	return nil
}

type ClockMock struct{}

func (c *ClockMock) Now() time.Time {
	return timeParse("2021-11-03 15:00:00 JST")
}

// timeParse は 2006-01-02 15:04:05 MST 形式の時間を変換する。失敗した場合は例外
func timeParse(str string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05 MST", str)
	if err != nil {
		panic(err)
	}
	return t
}

func TestApplication_CheckDeletableTweet(t *testing.T) {
	clock := &Clock{}
	type Want struct {
		isOK bool
		err  error
	}
	for _, tt := range []struct {
		adapter AdapterMock
		tweet   Tweet
		name    string
		want    Want
	}{
		{
			adapter: AdapterMock{[]LikeData{}},
			tweet: Tweet{
				CreatedAt: timeParse("2021-11-01 15:00:00 JST"),
				ID:        "1",
				Text:      "テストツイート",
			},
			name: "24時間以上経過していて、自分がいいねしていないツイートは削除されること",
			want: Want{
				isOK: true,
				err:  nil,
			},
		},
		{
			adapter: AdapterMock{[]LikeData{}},
			tweet: Tweet{
				CreatedAt: timeParse("2021-11-03 15:00:00 JST"),
				ID:        "1",
				Text:      "テストツイート",
			},
			name: "24時間以上経過していないツイートは削除されないこと",
			want: Want{
				isOK: false,
				err:  &CheckFailedError{Message: "24時間経っていないツイートです"},
			},
		},
		{
			adapter: AdapterMock{[]LikeData{
				{
					ID:       "101",
					Name:     "me",
					Username: "iam",
				},
			}},
			tweet: Tweet{
				CreatedAt: timeParse("2021-11-01 15:00:00 JST"),
				ID:        "1",
				Text:      "テストツイート",
			},
			name: "24時間以上経過していても、自分がいいねしているツイートは削除されないこと",
			want: Want{
				isOK: false,
				err:  &CheckFailedError{Message: "あなたがいいねしたツイートです"},
			},
		},
		{
			adapter: AdapterMock{[]LikeData{
				{
					ID:       "101",
					Name:     "me",
					Username: "iam",
				},
			}},
			tweet: Tweet{
				CreatedAt: timeParse("2021-11-03 15:00:00 JST"),
				ID:        "1",
				Text:      "テストツイート",
			},
			name: "24時間以上経過していないツイートは、自分がいいねしたとしても削除されないこと",
			want: Want{
				isOK: false,
				err:  &CheckFailedError{Message: "24時間経っていないツイートです"},
			},
		},
		{
			adapter: AdapterMock{[]LikeData{
				{
					ID:       "102",
					Name:     "other",
					Username: "notme",
				},
			}},
			tweet: Tweet{
				CreatedAt: timeParse("2021-11-01 15:00:00 JST"),
				ID:        "1",
				Text:      "テストツイート",
			},
			name: "24時間以上経過していて、他人がいいねしているツイートは削除されること",
			want: Want{
				isOK: true,
				err:  nil,
			},
		},
		{
			adapter: AdapterMock{[]LikeData{
				{
					ID:       "102",
					Name:     "other",
					Username: "notme",
				},
			}},
			tweet: Tweet{
				CreatedAt: timeParse("2021-11-03 15:00:00 JST"),
				ID:        "1",
				Text:      "テストツイート",
			},
			name: "24時間以上経過していない、他人がいいねしているツイートは削除されないこと",
			want: Want{
				isOK: false,
				err:  &CheckFailedError{Message: "24時間経っていないツイートです"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			app := Application{
				Client:      &tt.adapter,
				MyTwitterID: "101",
				clock:       clock,
			}
			isOK, err := app.CheckDeletableTweet(tt.tweet)
			assert.Exactly(t, tt.want.isOK, isOK)
			if diff := cmp.Diff(tt.want.err, err); diff != "" {
				t.Errorf("Error value is mismatch (-want +actual):\n%s", diff)
			}
		})
	}
}
