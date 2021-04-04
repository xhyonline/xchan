package server

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xchan/mod"
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
	// 先判断有没有设置七牛?
	isSet, err := h.s.ExistsSettingStoryType(mod.StoreType.QiNiu)
	if err != nil {
		c.JSON(200, Response(400, "获取是否设置存储类型失败,错误:"+err.Error(), nil))
		return
	}
	// 如果设置了,初始化七牛管理器
	if isSet {
		err = h.s.SetOSSManager()
		if err != nil {
			c.JSON(200, Response(400, err.Error(), nil))
			return
		}
	}
	// 在判断是否设置了本地存储
	isSet, err = h.s.ExistsSettingStoryType(mod.StoreType.Local)
	if err != nil {
		c.JSON(200, Response(400, "获取是否设置存储类型失败,错误:"+err.Error(), nil))
		return
	}
	// 如果设置了,初始化本地路径
	if isSet {
		err = h.s.SetLocalStorePath()
		if err != nil {
			c.JSON(200, Response(400, err.Error(), nil))
			return
		}

	}
	log.Info("路径:", h.s.Path)
	// 最后获取存储类型,并设置
	h.s.StoreType, err = h.s.GetStoryType()
	if err != nil {
		c.JSON(200, Response(400, "获取存储类型失败"+err.Error(), nil))
		return
	}

	c.JSON(200, Response(200, "登录成功", nil))
}
