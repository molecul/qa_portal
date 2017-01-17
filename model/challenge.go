package model

type Challenge struct {
	ID           int64  `xorm:"pk autoincr"`
	Name         string `xorm:"varchar(64) not null"`
	InternalName string `xorm:"varchar(64) unique not null"`
	Image        string `xorm:"varchar(64)"` // Docker image name
	TargetPath   string `xorm:"varchar(64)"` // Where file with answer been stored in container
	Cmd          string `xorm:"varchar(64)"` // Command to check answer. Can be null
	Description  string
	Points       int64
	FromDB       bool `xorm:"-"`
}
