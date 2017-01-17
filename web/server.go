package web

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
)

type Configuration struct {
	Hostname  string
	UseHTTP   bool
	UseHTTPS  bool
	HTTPPort  int
	HTTPSPort int
	CertFile  string
	KeyFile   string
}

func (c *Configuration) getServerPort(isTls bool) int {
	if isTls {
		return c.HTTPSPort
	} else {
		return c.HTTPPort
	}
}

func (c *Configuration) getServerAddr(isTls bool) string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.getServerPort(isTls))
}

func (c *Configuration) runServer(handlers http.Handler, isTls bool) (err error) {
	addr := c.getServerAddr(isTls)
	logrus.Info("Listening on ", addr)
	if isTls {
		err = http.ListenAndServeTLS(addr, c.CertFile, c.KeyFile, handlers)
	} else {
		err = http.ListenAndServe(addr, handlers)
	}
	if err != nil {
		logrus.Fatalf("Web listener on %s failed with error: %v", addr, err)
	} else {
		logrus.Warningf("Web listener on %s stopped", addr)
	}
	return
}

func (cfg *Configuration) initRoutes(r *gin.Engine) {

}

func Run(cfg *Configuration) {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	cfg.initRoutes(r)

	logrus.Info("Starting web server")

	if cfg.UseHTTP && cfg.UseHTTPS {
		w := sync.WaitGroup{}
		w.Add(2)
		go func() {
			cfg.runServer(r, true)
			w.Done()
		}()
		go func() {
			cfg.runServer(r, false)
			w.Done()
		}()
		w.Wait()
	} else {
		cfg.runServer(r, cfg.UseHTTPS)
	}
	logrus.Print("Web server stopped")
}
