package domain

import "time"

type Tweet struct {
	CreatedAt time.Time
	ID        string
	Text      string
}
type LikeData struct {
	ID       string
	Name     string
	Username string
}
