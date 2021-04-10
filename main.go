package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xchan/middleware"
	"github.com/xhyonline/xchan/server"
	"github.com/xhyonline/xutil/xlog"
	"net/http"
)

var log = xlog.Get(true)

func main() {
	log.Infof("程序启动了")
	s := server.GetService()
	h := server.NewHandler(s)
	g := gin.Default()
	// 存储最大限制
	g.MaxMultipartMemory = 20480 << 20
	// 修改模板标签
	g.Delims("<go", "go>")
	// 前端 HTML 文件
	g.LoadHTMLGlob("./views/layui/views/*")

	// css 、 js 等静态资源文件
	g.StaticFS("/layuiadmin", http.Dir("./views/layui/layuiadmin"))
	// jquery 拖拽上传插件
	g.StaticFS("/drop", http.Dir("./views/dist"))
	// 本地文件存储位置
	g.StaticFS("/file-save-dir", http.Dir("./file-save-dir"))

	// 前端路由组与中间件
	front := g.Group("")
	// 检查是否已安装、是否已登录、设置存储类型
	front.Use(middleware.CheckInstall, middleware.HaveLogin)
	{
		front.GET("/", h.Login)
		// 登录
		front.POST("/login-check", h.LoginCheck)
	}
	// 安装路由
	g.POST("/install", h.Install)
	//
	g.GET("/install-view", h.InstallView)
	// 后台路由组
	admin := g.Group("/admin")
	admin.Use(middleware.Auth, middleware.CheckInstall)
	{
		// 后台首页
		admin.GET("/index", h.Admin)
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
		// 执行新增用户
		admin.POST("/exec/add-user", h.AddUserItem)
		// 执行修改用户
		admin.POST("/exec/update-user", h.UpdateUserItem)
		// 设置界面
		admin.GET("/setting", h.Setting)
	}

	err := g.Run("0.0.0.0:80")
	if err != nil {
		panic(err)
	}
}
