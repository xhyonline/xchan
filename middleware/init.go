package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xchan/server"
)

// InitInstall 初始化,查看用户是否安装了本软件
func InitInstall(c *gin.Context) {
	s := server.GetService()
	install, err := s.Config.GetSection("install")
	if err != nil {
		c.JSON(200, server.Response(400, "初始化失败"+err.Error(), nil))
		c.Abort()
		return
	}
	// 如果没安装.直接弹出安装界面
	if install["have_install"] == "false" {
		log.Infof("用户还没安装")
		c.HTML(200, "install.html", nil)
		c.Abort()
		return
	}
	// 如果安装了,看看 ping 的通不
	err = s.DB.DB().Ping()
	if err != nil {
		c.JSON(200, server.Response(400, "数据库链接失败"+err.Error(), nil))
		c.Abort()
		return
	}
	// 如果 ping 的通
	c.Next()
}
