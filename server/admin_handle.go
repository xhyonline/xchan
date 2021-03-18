package server

import "github.com/gin-gonic/gin"

// Admin 后台首页
func (h *Handler) Admin(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

// Console 后台首页
func (h *Handler) Console(c *gin.Context) {
	c.HTML(200, "console.html", nil)
}
