package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

type Configuration struct {
	Driver string
	DSN    string
}

var engine *xorm.Engine

func Init(cfg *Configuration) (err error) {
	if engine, err = xorm.NewEngine(cfg.Driver, cfg.DSN); err != nil {
		return
	}

	engine.SetLogger(&XormLogger{})

	return engine.Ping()
}

func Get() *xorm.Engine {
	return engine
}
