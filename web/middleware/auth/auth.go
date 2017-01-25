package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/molecul/qa_portal/model"
	"github.com/molecul/qa_portal/web/views"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// User is a retrieved and authenticated user.
type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Hd            string `json:"hd"`
}

func (gu *User) FillUserInfo(u *model.User) *model.User {
	u.Email = strings.ToLower(gu.Email)
	u.EmailVerified = gu.EmailVerified
	u.Name = gu.Name
	u.Picture = gu.Picture
	return u
}

var conf *oauth2.Config
var userIdKey = "model.user.id"
var userKey = "model.user"

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// Setup the authorization path
func Setup(redirectURL, clientId, clientSecret string, scopes []string) {
	conf = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}

func Login(ctx *gin.Context, redirect string) {
	state := randToken()
	session := sessions.Default(ctx)
	session.Set("state", state)
	session.Set("redirect", redirect)
	session.Save()

	views.RenderRedirect(ctx, conf.AuthCodeURL(state))
}

func Logout(ctx *gin.Context, redirect string) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	views.RenderRedirect(ctx, redirect)
}

func userDataUpdate(ctx *gin.Context, gu *User) {
	usr, err := model.GetUserByEmail(gu.Email)
	if err != nil {
		views.RenderError(ctx, err)
		return
	}
	if usr == nil {
		if err = model.CreateUser(gu.FillUserInfo(new(model.User))); err != nil {
			views.RenderError(ctx, err)
			return
		}
	} else {
		if err = gu.FillUserInfo(usr).Update(); err != nil {
			views.RenderError(ctx, err)
			return
		}
	}
	session := sessions.Default(ctx)
	session.Set(userIdKey, usr.ID)
	session.Save()
}

func AuthRedirectHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Handle the exchange code to initiate a transport.
		session := sessions.Default(ctx)
		retrievedState := session.Get("state")
		if retrievedState != ctx.Query("state") {
			ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
			return
		}

		tok, err := conf.Exchange(oauth2.NoContext, ctx.Query("code"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		client := conf.Client(oauth2.NoContext, tok)
		email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer email.Body.Close()
		data, err := ioutil.ReadAll(email.Body)
		if err != nil {
			glog.Errorf("[Gin-OAuth] Could not read Body: %s", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		user := new(User)
		err = json.Unmarshal(data, user)
		if err != nil {
			glog.Errorf("[Gin-OAuth] Unmarshal userinfo failed: %s", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// save userinfo
		userDataUpdate(ctx, user)
		if redirect := session.Get("redirect"); redirect != nil {
			views.RenderRedirect(ctx, redirect.(string))
		} else {
			views.RenderRedirect(ctx, "/")
		}
	}
}

func GetUser(ctx *gin.Context) *model.User {
	if u, ex := ctx.Get(userKey); ex {
		return u.(*model.User)
	} else {
		return nil
	}
}

func LoginRequired(h gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ex := ctx.Get(userKey); !ex {
			Login(ctx, ctx.Request.URL.String())
		} else {
			h(ctx)
		}
	}
}

func UserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userId := session.Get(userIdKey)
		if userId != nil {
			usr, err := model.GetUserById(userId.(int64))
			if err != nil {
				panic(err)
			}
			ctx.Set(userKey, usr)
		}
	}
}
