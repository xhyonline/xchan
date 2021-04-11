package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xchan/server"
)

// CheckInstall 初始化,查看用户是否安装了本软件
func CheckInstall(c *gin.Context) {
	s := server.GetService()
	install, err := s.Config.GetSection("install")
	if err != nil {
		c.JSON(200, server.Response(400, "初始化失败"+err.Error(), nil))
		c.Abort()
		return
	}
	// 如果没安装.直接弹出安装界面
	if install["have_install"] == "false" {
		c.SetCookie("token", "", -1, "/", c.Request.Host, false, true)
		c.Redirect(307, "/install-view")
		c.Abort()
		return
	}
	// 如果安装了,再检查一次
	if err := s.CheckInstalled(); err != nil {
		c.JSON(200, server.Response(400, "安装后检查失败"+err.Error(), nil))
		c.Abort()
		return
	}
	c.Next()
}
