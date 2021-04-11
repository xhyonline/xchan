package server

import (
	"github.com/gin-gonic/gin"
)

type User struct {
	UserName    string `json:"user_name" form:"user_name"`
	Password    string `json:"password" form:"password"`
	RePassword  string `json:"re_password" form:"re_password"`
	OldPassword string `json:"old_password" form:"old_password"`
	NewPassword string `json:"new_password" form:"new_password"`
}

// AddUserItem 新增一个用户
func (h *Handler) AddUserItem(c *gin.Context) {
	item := new(User)
	err := c.Bind(item)
	if err != nil {
		c.JSON(200, Response(400, err.Error(), nil))
		return
	}
	if item.Password != item.RePassword {
		c.JSON(200, Response(400, "两次密码不匹配", nil))
		return
	}
	err = h.s.AddUser(item)
	if err != nil {
		c.JSON(200, Response(400, err.Error(), nil))
		return
	}
	c.JSON(200, Response(200, "新增成功", nil))
}

// UpdateUserItem 修改一个用户
func (h *Handler) UpdateUserItem(c *gin.Context) {
	item := new(User)
	err := c.Bind(item)
	if err != nil {
		c.JSON(200, Response(400, err.Error(), nil))
		return
	}
	err = h.s.UpdateUser(item)
	if err != nil {
		c.JSON(200, Response(400, err.Error(), nil))
		return
	}
	c.SetCookie("token", "", -1, "/", c.Request.Host, false, true)
	c.JSON(200, Response(200, "修改成功", nil))
	return
}
