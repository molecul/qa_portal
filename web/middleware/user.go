package middleware

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/molecul/qa_portal/model"
)

var userIdKey = "model.user.id"
var userKey = "model.user"

func User() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get(userIdKey)
		if userId != nil {
			usr, err := model.GetUserById(userId.(int64))
			if err != nil {
				panic(err)
			}
			c.Set(userKey, usr)
		}
	}
}

func UserMust(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ex := c.Get(userKey); !ex {
			c.Redirect(http.StatusMovedPermanently, "/login")
		} else {
			h(c)
		}
	}
}

func UserFromContext(c *gin.Context) *model.User {
	if u, ex := c.Get(userKey); ex {
		return u.(*model.User)
	} else {
		return nil
	}
}

func UserSessionSet(c *gin.Context, userid int64) {
	session := sessions.Default(c)
	session.Set(userIdKey, userid)
	session.Save()
}
