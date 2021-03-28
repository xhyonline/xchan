package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xchan/server"
)

// Auth 初步鉴权
func Auth(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		fmt.Println("权限不足1")
		c.Redirect(301, "/")
		return
	}
	s := server.GetService()
	if _, err = s.ParseToken(token); err != nil {
		fmt.Println("权限不足2")
		c.Redirect(301, "/")
		return
	}
	c.Next()
}
