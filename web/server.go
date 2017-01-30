package web

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	eztemplate "github.com/michelloworld/ez-gin-template"
	"github.com/molecul/qa_portal/util/isdebug"
	"github.com/molecul/qa_portal/web/handlers"
	"github.com/molecul/qa_portal/web/middleware/auth"
)

var googleScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
	// You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
}

type GoogleOAuthConfig struct {
	Secret          string
	SessionName     string
	OAuthClientId   string
	OAuthSecret     string
	SessionDuration time.Duration
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
	auth.Setup(cfg.getHostname()+"/auth/", cfg.GoogleOAuth.OAuthClientId, cfg.GoogleOAuth.OAuthSecret,
		googleScopes, cfg.GoogleOAuth.SessionDuration*time.Second)
}

func (cfg *Configuration) initRoutes(r *gin.Engine) {
	render := eztemplate.New()
	render.TemplatesDir = "web/templates/"
	render.Layout = "base"
	render.TemplateFuncMap = initTemplatesFuncs()
	if isdebug.Is {
		render.Debug = true
	}
	r.HTMLRender = render.Init()

	r.StaticFS("/static", http.Dir("./web/static"))

	r.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte(cfg.GoogleOAuth.Secret))))
	r.Use(auth.UserMiddleware())

	r.GET("/", handlers.MainPageHandler)

	// Auth section
	r.GET("/auth", auth.AuthRedirectHandler())
	r.GET("/login", handlers.LoginHandler)
	r.GET("/logout", handlers.LogoutHandler)
	r.GET("/challenges", handlers.ChallengesWebHandler)
	r.GET("/scoreboard", handlers.ScoreboardHandler)
	r.GET("/profile", auth.LoginRequired(handlers.ProfileHandler))
	r.GET("/solve/:challenge", auth.LoginRequired(handlers.SolveHandlerGet))
	r.POST("/solve/:challenge", auth.LoginRequired(handlers.SolveHandlerPost))
	r.GET("/test/:test", auth.LoginRequired(handlers.TestHandler))

	// Api section
	api := r.Group("/api")
	api.GET("/healthcheck", auth.LoginRequired(handlers.DockerHealthCheckHandler))
	api.GET("/userinfo", auth.LoginRequired(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"user": auth.GetUser(ctx)})
	}))
	api.GET("/users", auth.LoginRequired(handlers.UsersHandler))
	api.GET("/history", auth.LoginRequired(handlers.ProfileHistoryHandler))
	api.GET("/challenges", handlers.ChallengesHandler)
}

func Run(cfg *Configuration) {
	if !isdebug.Is {
		gin.SetMode(gin.ReleaseMode)
	}
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
