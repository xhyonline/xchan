package server

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xchan/mod"
	"strconv"
)

// UploadQiNiu 上传接口
func (h *Handler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, Response(400, "上传失败"+err.Error(), nil))
		return
	}
	token, err := c.Cookie("token")
	t, err := h.s.ParseToken(token)
	if err != nil {
		c.JSON(200, Response(400, "上传失败"+err.Error(), nil))
		return
	}
	var src string
	switch h.s.StoreType {
	case mod.StoreType.QiNiu: // 如果存储类型是七牛
		src, err = h.s.UploadQiNiu(file, t.Username)
	case mod.StoreType.Local: // 如果存储类型是本地
		src, err = h.s.UploadLocal(file, t.Username)
	}
	if err != nil {
		c.JSON(200, Response(400, "上传失败"+err.Error(), nil))
		return
	}
	c.JSON(200, Response(200, "上传成功", gin.H{
		"src": src,
	}))
}

// Remove 删除接口
func (h *Handler) Remove(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(200, Response(400, "移除失败,请输入正确的ID", nil))
		return
	}
	err := h.s.Remove(id)
	if err != nil {
		c.JSON(200, Response(400, "移除失败"+err.Error(), nil))
		return
	}
	c.JSON(200, Response(200, "移除成功", nil))
}

type RemoveList struct {
	List []string `json:"list"`
}

// Remove 删除接口
func (h *Handler) RemoveList(c *gin.Context) {
	form := new(RemoveList)
	err := c.Bind(form)
	if err != nil {
		c.JSON(200, Response(400, "移除失败,请输入正确的ID列表", nil))
	}
	err = h.s.RemoveList(form.List)
	if err != nil {
		c.JSON(200, Response(400, "移除失败"+err.Error(), nil))
		return
	}
	c.JSON(200, Response(200, "移除成功", nil))
	return
}

// GetList 获取文件列表
func (h *Handler) GetList(c *gin.Context) {
	p := c.Query("page")
	l := c.Query("limit")
	page, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(200, Response(400, err.Error(), nil))
		return
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		c.JSON(200, Response(400, err.Error(), nil))
		return
	}

	resp, total, err := h.s.GetListByPaging(page, limit)
	if err != nil {
		c.JSON(200, Response(400, err.Error(), nil))
		return
	}
	m := Response(0, "获取成功", resp)
	m["count"] = total
	c.JSON(200, m)

}
