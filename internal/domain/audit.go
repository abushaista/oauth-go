package domain

import "time"

type Audit struct {
	ID        string
	UserID    string
	ClientID  string
	Action    string
	Details   string
	IPAddress string
	CreatedAt time.Time
}
