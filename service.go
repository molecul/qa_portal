package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/molecul/qa_portal/checker"
	"github.com/molecul/qa_portal/model"
	"github.com/molecul/qa_portal/util/database"
	"github.com/molecul/qa_portal/web"
	"gopkg.in/yaml.v2"
)

const TimestampFormat = "06-01-02 15:04:05.000"

type PrettyFormatter struct{}

func (f *PrettyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s|%7s|%s\n", entry.Time.Format(TimestampFormat), entry.Level.String(), entry.Message)), nil
}

type Configuration struct {
	Web      *web.Configuration
	Model    *model.Configuration
	Checker  *checker.Configuration
	Database *database.Configuration
}

func init() {
	logrus.SetFormatter(&PrettyFormatter{})

	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func Start(c *Configuration) {
	logrus.Infof("Starting server")

	if err := database.Init(c.Database); err != nil {
		logrus.Fatalf("Error when connecting to db: %v", err)
	}

	if err := model.Init(c.Model); err != nil {
		logrus.Fatalf("Error syncing db: %v", err)
	}

	if err := checker.Init(c.Checker); err != nil {
		logrus.Fatalf("Error when creating checker: %v", err)
	}

	web.Run(c.Web)
}

func main() {
	var confPath string
	flag.StringVar(&confPath, "config", "config.yml", "path to config yaml file")
	flag.Parse()

	cfgRaw, err := ioutil.ReadFile(confPath)
	if err != nil {
		logrus.Fatalf("Cannot open config file %s: %v", confPath, err)
	}
	cfg := new(Configuration)
	if err := yaml.Unmarshal(cfgRaw, cfg); err != nil {
		logrus.Fatalf("Cannot unmarshal config: %v", err)
	}
	Start(cfg)
}
