package jsonapi

import (
	"time"
)

type Touched struct {
	Created  TouchedIdentity `json:"created"`
	Modified TouchedIdentity `json:"modified"`
}

type TouchedIdentity struct {
	UserID string    `json:"user_id"`
	Time   time.Time `json:"time"`
}
