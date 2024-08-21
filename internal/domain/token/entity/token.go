package entity

import "time"

type AuthToken struct {
	TokenType   string
	TokenValue  string
	ExpiresTime *time.Time
}
