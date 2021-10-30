package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

func UserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID := session.Get("id")

		if sessionID == nil {
			session.Set("id", uuid.NewV4())
			session.Save()
		}
	}
}
