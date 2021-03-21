package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xchan/middleware"
	"github.com/xhyonline/xchan/mod"
	"github.com/xhyonline/xchan/server"
	"net/http"
)

func main() {
	s := server.GetService()
	s.DB.AutoMigrate(&mod.User{}, &mod.OSS{})
	h := server.NewHandler(s)
	g := gin.Default()
	g.MaxMultipartMemory = 20480 << 20 // 8 MiB
	// 前端 HTML 文件
	g.LoadHTMLGlob("./views/layui/views/*")
	// css 、 js 等静态资源文件
	g.StaticFS("/layuiadmin", http.Dir("./views/layui/layuiadmin"))

	// 登录
	g.GET("/", h.Login)
	g.POST("/login-check", h.LoginCheck)

	// 后台路由组
	admin := g.Group("/admin")
	{
		admin.Use(middleware.Auth)
		// 后台首页
		admin.GET("/", h.Admin)
		// 控制台
		admin.GET("/console", h.Console)
		// 文件上传接口
		admin.POST("/upload", h.Upload)
		// 删除接口
		admin.GET("/remove", h.Remove)
		// 文件列表
		admin.GET("/list")
	}

	err := g.Run("0.0.0.0:80")
	if err != nil {
		panic(err)
	}
}
