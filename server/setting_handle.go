package server

import "github.com/gin-gonic/gin"

func (h *Handler) Setting(c *gin.Context) {
	c.HTML(200, "setting.html", nil)
}
