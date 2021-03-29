package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/xhyonline/xchan/server"
)

// Auth 初步鉴权
func Auth(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Infof("请求权限不足")
		c.Redirect(307, "/")
		return
	}
	s := server.GetService()
	if _, err = s.ParseToken(token); err != nil {
		log.Infof("请求权限不足")
		c.Redirect(307, "/")
		return
	}
	c.Next()
}
