package middleware

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

func UserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID := session.Get("id")

		session.Options(sessions.Options{
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			MaxAge:   0,
		})

		if sessionID == nil {
			session.Set("id", uuid.NewV4())
		}

		session.Save()

		log.Default().Printf("UserID %s accessing %s", session.Get("id"), c.Request.URL)
	}
}
