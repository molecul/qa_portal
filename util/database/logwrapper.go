package database

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-xorm/core"
)

// XormLogger implement go-xorm/core.ILogger
type XormLogger struct {
	showSQL bool
}

func (l *XormLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

func (s *XormLogger) IsShowSQL() bool {
	return s.showSQL
}

func (s *XormLogger) Level() core.LogLevel {
	return core.LOG_DEBUG
}

func (s *XormLogger) SetLevel(core.LogLevel) {}

func (s *XormLogger) Debug(v ...interface{}) {
	logrus.Debug(v)
}

func (s *XormLogger) Debugf(format string, v ...interface{}) {
	logrus.Debugf(format, v)
}

func (s *XormLogger) Error(v ...interface{}) {
	logrus.Error(v)
}

func (s *XormLogger) Errorf(format string, v ...interface{}) {
	logrus.Errorf(format, v)
}

func (s *XormLogger) Info(v ...interface{}) {
	logrus.Info(v)
}

func (s *XormLogger) Infof(format string, v ...interface{}) {
	logrus.Infof(format, v)
}

func (s *XormLogger) Warn(v ...interface{}) {
	logrus.Warn(v)
}

func (s *XormLogger) Warnf(format string, v ...interface{}) {
	logrus.Warnf(format, v)
}
