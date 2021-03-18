package server

import (
	"github.com/gin-gonic/gin"
	"time"
)

// Login 登录时
func (h *Handler) Login(c *gin.Context) {
	token, err := c.Cookie("token")
	if err == nil {
		// cookie 没过期
		if _, err = h.s.ParseToken(token); err == nil {
			c.Redirect(301, "/admin")
			return
		}
	}
	c.HTML(200, "login.html", nil)
}

// LoginCheck 校验基本信息
func (h *Handler) LoginCheck(c *gin.Context) {
	id := c.PostForm("username")
	password := c.PostForm("password")
	isTrue, err := h.s.CheckLogin(id, password)
	if err != nil {
		c.JSON(400, response(400, err.Error(), nil))
		return
	}
	if isTrue {
		// 生成 token
		token, err := h.s.GenerateToken(id)
		if err != nil {
			c.JSON(400, response(400, err.Error(), nil))
			return
		}
		c.SetCookie("token", token, int(time.Hour*72), "/", "localhost", false, true)
		c.JSON(200, response(200, "登录成功", nil))
		return
	}
	c.JSON(200, response(400, "验证失败", nil))
}
