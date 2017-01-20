package web

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/molecul/qa_portal/web/handlers"
	"github.com/zalando/gin-oauth2/google"
)

var googleScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	// You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
}

type GoogleOAuthConfig struct {
	Secret      string
	SessionName string
	CredFile    string
}

type Configuration struct {
	Hostname    string
	Listen      string
	UseHTTP     bool
	UseHTTPS    bool
	HTTPPort    int
	HTTPSPort   int
	CertFile    string
	KeyFile     string
	CredFile    string
	GoogleOAuth GoogleOAuthConfig
}

func (c *Configuration) getServerPort(isTls bool) int {
	if isTls {
		return c.HTTPSPort
	} else {
		return c.HTTPPort
	}
}

func (c *Configuration) getServerAddr(isTls bool) string {
	return fmt.Sprintf("%s:%d", c.Listen, c.getServerPort(isTls))
}

func (c *Configuration) getHostname() string {
	if c.UseHTTPS {
		return fmt.Sprintf("https://%s:%d", c.Hostname, c.HTTPSPort)
	} else {
		return fmt.Sprintf("http://%s:%d", c.Hostname, c.HTTPPort)
	}
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

func (cfg *Configuration) setupGoogle() {
	gc := &cfg.GoogleOAuth
	google.Setup(cfg.getHostname()+"/auth/", gc.CredFile, googleScopes, []byte(gc.Secret))
}

func (cfg *Configuration) initRoutes(r *gin.Engine) {
	r.StaticFS("/static", http.Dir("./web/static"))

	r.LoadHTMLGlob("./web/templates/*")

	r.Use(google.Session(cfg.GoogleOAuth.SessionName))
	r.GET("/login", google.LoginHandler)

	// protected url group
	private := r.Group("/auth")
	private.Use(google.Auth())
	{
		private.GET("/", webHandlers.UserInfoHandler)
		private.GET("/api", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "Hello from private for groups"})
		})
	}

	r.GET("/", webHandlers.MainPageHandler)

}

func Run(cfg *Configuration) {
	r := gin.Default()

	cfg.setupGoogle()
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
