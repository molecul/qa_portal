package middleware

import (
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

func UserFromContext(c *gin.Context) *model.User {
	return c.MustGet(userKey).(*model.User)
}

func UserSessionSet(c *gin.Context, userid int64) {
	sessions.Default(c).Set(userIdKey, userid)
}
