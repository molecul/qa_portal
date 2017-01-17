package model

import (
	"time"
)

type Test struct {
	ID          int64 `xorm:"pk autoincr"`
	ChallengeID int64
	UserID      int64
	IsSucess    bool
	InputFile   string
	OutputFile  string
	Created     time.Time `xorm:"created"`
	Checked     time.Time
	Duration    time.Duration
}
