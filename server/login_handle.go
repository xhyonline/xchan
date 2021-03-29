package server

import (
	"github.com/gin-gonic/gin"
	"time"
)

// Login 登录时
func (h *Handler) Login(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.HTML(200, "login.html", nil)
		return
	}
	if _, err = h.s.ParseToken(token); err != nil {
		c.HTML(200, "login.html", nil)
		return
	}
	// 否则从登录界面直接重定向到
	c.Redirect(307, "/admin/index")
}

// LoginCheck 校验基本信息
func (h *Handler) LoginCheck(c *gin.Context) {
	id := c.PostForm("username")
	password := c.PostForm("password")
	isTrue, err := h.s.CheckIDAndPass(id, password)
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
		c.SetCookie("token", token, int(time.Hour*72), "/", c.Request.Host, false, true)
		c.JSON(200, response(200, "登录成功", nil))
		return
	}
	c.JSON(200, response(400, "验证失败", nil))
}
