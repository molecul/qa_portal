package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/molecul/qa_portal/util/database"
)

type Test struct {
	Id          int64      `xorm:"pk autoincr"`
	ChallengeId int64      `xorm:"notnull"`
	UserId      int64      `xorm:"notnull"`
	IsSucess    bool       `xorm:"notnull"`
	Created     time.Time  `xorm:"notnull created"`
	Checked     *time.Time `xorm:"null"`
	Duration    time.Duration
}

func GetTestById(id int64) (*Test, error) {
	t := new(Test)
	has, err := database.Get().Id(id).Get(t)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return t, nil
}

func CreateTest(t *Test, input []byte) error {
	if _, err := database.Get().InsertOne(t); err != nil {
		return err
	}

	inputFile := t.GetInputFileName()
	if err := os.MkdirAll(filepath.Dir(inputFile), 0755); err != nil {
		return fmt.Errorf("Error creating task dirs: %v", err)
	}
	if err := ioutil.WriteFile(inputFile, input, 0755); err != nil {
		return fmt.Errorf("Error creating task file: %v", err)
	}

	return nil
}

func (t *Test) getBaseFileName(suffix string) string {
	return filepath.Join(configuration.LocalTestFiles, strconv.FormatInt(t.ChallengeId, 16), strconv.FormatInt(t.Id, 16)+suffix)
}

func (t *Test) GetInputFileName() string {
	return t.getBaseFileName(".input")
}

func (t *Test) GetOutputFileName() string {
	return t.getBaseFileName(".output")
}

func (t *Test) Update(output []byte) error {
	_, err := database.Get().Id(t.Id).AllCols().Update(t)
	if output != nil {
		err := ioutil.WriteFile(t.GetOutputFileName(), output, 0755)
		if err != nil {
			return err
		}
	}
	return err
}

func Tests(page, pageSize int, onlyUntested bool, forUser int64) ([]*Test, error) {
	tests := make([]*Test, 0, pageSize)
	q := database.Get().NewSession()
	if pageSize != 0 {
		if page == 0 {
			q = q.Limit(pageSize)
		} else {
			q = q.Limit(pageSize, (page-1)*pageSize)
		}
	}
	if onlyUntested {
		q = q.Where("checked is null")
	}
	if forUser != 0 {
		q = q.Where("user_id == ?", forUser)
	}
	return tests, q.Find(&tests)
}
