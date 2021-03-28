package server

import (
	"github.com/gin-gonic/gin"
)

// Admin 后台首页
func (h *Handler) Admin(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

// Console 后台首页
func (h *Handler) Console(c *gin.Context) {
	c.HTML(200, "console.html", nil)
}

// List 文件列表
func (h *Handler) List(c *gin.Context) {
	c.HTML(200, "list.html", nil)
}

// AddUser 新增用户
func (h *Handler) AddUser(c *gin.Context) {
	c.HTML(200, "add-user.html", nil)
}

// UpdateUser 新增用户
func (h *Handler) UpdateUser(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		panic(err)
	}
	clm, err := h.s.ParseToken(token)
	if err != nil {
		panic(err)
	}
	c.HTML(200, "update-user.html", gin.H{
		"username": clm.Username,
	})
}
