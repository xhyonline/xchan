package server

import (
	"github.com/gin-gonic/gin"
	"time"
)

// Login 登录时
func (h *Handler) Login(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

// LoginCheck 校验基本信息
func (h *Handler) LoginCheck(c *gin.Context) {
	id := c.PostForm("username")
	password := c.PostForm("password")
	isTrue, err := h.s.CheckIDAndPass(id, password)
	if err != nil {
		c.JSON(200, Response(400, err.Error(), nil))
		return
	}
	if !isTrue {
		c.JSON(200, Response(400, "验证失败", nil))
		return
	}
	// 生成 token
	token, err := h.s.GenerateToken(id)
	if err != nil {
		c.JSON(200, Response(400, err.Error(), nil))
		return
	}
	c.SetCookie("token", token, int(time.Hour*72), "/", c.Request.Host, false, true)
	c.JSON(200, Response(200, "登录成功", nil))
}
