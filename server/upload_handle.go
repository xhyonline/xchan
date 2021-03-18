package server

import "github.com/gin-gonic/gin"

// Upload 上传接口
func (h *Handler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, response(400, "上传失败"+err.Error(), nil))
		return
	}
	src, err := h.s.Upload(file)
	if err != nil {
		c.JSON(200, response(400, "上传失败"+err.Error(), nil))
		return
	}
	c.JSON(200, response(200, "上传成功", gin.H{
		"src": src,
	}))
}
