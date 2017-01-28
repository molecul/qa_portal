package model

import (
	"time"

	"github.com/molecul/qa_portal/util/database"
)

type Challenge struct {
	Id           int64  `xorm:"pk autoincr"`
	Name         string `xorm:"not null"`
	InternalName string `xorm:"unique not null"`
	Image        string `json:"-"` // Docker image name
	TargetPath   string `json:"-"` // Where file with been stored in container
	Cmd          string `json:"-"` // Command to check answer. Can be null
	Inject       string `json:"-"`
	Description  string `json:"-"`
	Points       int64
	Created      time.Time `xorm:"created"`
}

func GetChallengeById(id int64) (*Challenge, error) {
	c := new(Challenge)
	has, err := database.Get().Id(id).Get(c)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return c, nil
}

func GetChallengeByInternalName(name string) (*Challenge, error) {
	c := &Challenge{InternalName: name}
	has, err := database.Get().Get(c)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return c, nil
}

func CreateChallenge(c *Challenge) error {
	_, err := database.Get().InsertOne(c)
	return err
}

func (c *Challenge) Update() error {
	_, err := database.Get().Id(c.Id).AllCols().Update(c)
	return err
}

func (c *Challenge) IsEqual(o *Challenge) bool {
	return c.Name == o.Name &&
		c.Image == o.Image &&
		c.TargetPath == o.TargetPath &&
		c.Cmd == o.Cmd &&
		c.Description == o.Description &&
		c.Points == o.Points &&
		c.Inject == o.Inject
}

func (c *Challenge) UpdateWithInfoFrom(o *Challenge) error {
	c.Name = o.Name
	c.Image = o.Image
	c.TargetPath = o.TargetPath
	c.Cmd = o.Cmd
	c.Description = o.Description
	c.Points = o.Points
	c.Inject = o.Inject
	return c.Update()
}

func Challenges(page, pageSize int, order string) ([]*Challenge, error) {
	challenges := make([]*Challenge, 0, pageSize)

	if order == "asc" {
		return challenges, database.Get().Limit(pageSize, (page-1)*pageSize).Asc("id").Find(&challenges)
	} else {
		return challenges, database.Get().Limit(pageSize, (page-1)*pageSize).Desc("id").Find(&challenges)
	}
}
