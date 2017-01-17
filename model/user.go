package model

import (
	"time"
)

type User struct {
	ID      int64 `xorm:"pk autoincr"`
	Score   int64
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	// TODO OAuth sutff
}
