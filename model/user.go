package model

import (
	"time"
)

type User struct {
	ID        int64
	Score     int64
	CreatedAt time.Time
	UpdatedAt time.Time
	// TODO OAuth sutff
}
