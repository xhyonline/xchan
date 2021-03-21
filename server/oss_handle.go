package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Upload 上传接口
func (h *Handler) Upload(c *gin.Context) {
	fmt.Println("走到我这里了")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, response(400, "上传失败"+err.Error(), nil))
		return
	}

	token, err := c.Cookie("token")
	t, err := h.s.ParseToken(token)
	if err != nil {
		c.JSON(200, response(400, "上传失败"+err.Error(), nil))
		return
	}
	src, err := h.s.Upload(file, t.Username)
	if err != nil {
		c.JSON(200, response(400, "上传失败"+err.Error(), nil))
		return
	}
	c.JSON(200, response(200, "上传成功", gin.H{
		"src": src,
	}))
}

// Remove 删除接口
func (h *Handler) Remove(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(200, response(400, "移除失败,请输入正确的图片路径", nil))
		return
	}

	err := h.s.Remove(path)
	if err != nil {
		c.JSON(200, response(400, "移除失败"+err.Error(), nil))
		return
	}
	c.JSON(200, response(200, "移除成功", nil))
}

// List 文件列表
func (h *Handler) List(c *gin.Context) {

}
