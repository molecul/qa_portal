package model

import (
	"time"
)

type Test struct {
	ID          int64
	ChallengeID int64
	UserID      int64
	IsSucess    bool
	InputFile   string
	OutputFile  string
	CreatedAt   time.Time
	CheckedAt   time.Time
	Duration    time.Duration
}
