package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xchan/server"
)

// Auth 初步鉴权
func Auth(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.Redirect(301, "/")
		return
	}
	s := server.GetService()
	if _, err = s.ParseToken(token); err != nil {
		c.Redirect(301, "/")
		return
	}
	c.Next()
}
