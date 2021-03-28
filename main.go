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
	// 修改模板标签
	g.Delims("<go", "go>")
	// 前端 HTML 文件
	g.LoadHTMLGlob("./views/layui/views/*")

	// css 、 js 等静态资源文件
	g.StaticFS("/layuiadmin", http.Dir("./views/layui/layuiadmin"))
	// jquery 拖拽上传插件
	g.StaticFS("/drop", http.Dir("./views/dist"))

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
		// 批量删除
		admin.POST("/remove-list", h.RemoveList)
		// 文件列表
		admin.GET("/list", h.List)
		// 文件列表获取数据
		admin.GET("/list/get", h.GetList)
		// 新增用户界面
		admin.GET("/add-user", h.AddUser)
		// 修改用户信息界面
		admin.GET("/update-user", h.UpdateUser)

		admin.POST("/exec/add-user", h.AddUserItem)

		admin.POST("/exec/update-user", h.UpdateUserItem)
	}

	err := g.Run("0.0.0.0:80")
	if err != nil {
		panic(err)
	}
}
