package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xchan/server"
	"github.com/xhyonline/xutil/xlog"
)

var logger = xlog.Get().Debug()

// Auth 初步鉴权
func Auth(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Infof("请求后台权限不足")
		c.Redirect(307, "/")
		return
	}
	s := server.GetService()
	if _, err = s.ParseToken(token); err != nil {
		logger.Infof("请求后台权限不足")
		c.Redirect(307, "/")
		return
	}
	c.Next()
}

// HaveLogin 是否登录过了
func HaveLogin(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.Next()
		return
	}
	s := server.GetService()
	if _, err = s.ParseToken(token); err != nil {
		c.Next()
		return
	}
	// 跳转后台
	c.Redirect(307, "/admin/index")
}
