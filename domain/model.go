package domain

import "time"

type Tweet struct {
	CreatedAt       time.Time
	ID              string
	Text            string
	ConversationID  string
	ReferencedTweet []ReferencedTweet
}
type LikeData struct {
	ID       string
	Name     string
	Username string
}

type ReferencedTweet struct {
	Type string
	ID   string
}
