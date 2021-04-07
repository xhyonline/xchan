package server

import (
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/xhyonline/xchan/mod"
	"github.com/xhyonline/xutil/helper"
	"strings"
)

// Install 安装方法
func (h *Handler) Install(c *gin.Context) {
	host := c.PostForm("host")
	user := c.PostForm("username")
	password := c.PostForm("password")
	dbUser := c.PostForm("db_username")
	dbPassword := c.PostForm("db_password")
	port := c.PostForm("port")
	position := c.PostForm("position")
	key := c.PostForm("qiniu_key")
	secret := c.PostForm("qiniu_secret")
	bucket := c.PostForm("qiniu_bucket")
	qiNiuDomain := c.PostForm("qiniu_domain")
	localDomain := c.PostForm("local_domain")

	if position == "local" && !helper.IsURL(localDomain) {
		c.JSON(200, Response(400, "URL 不符合规范", nil))
		return
	}
	if position == "qiniu" && !helper.IsURL(qiNiuDomain) {
		c.JSON(200, Response(400, "URL 不符合规范", nil))
		return
	}

	// 尝试链接数据库
	err := h.s.ConnectDB(host, dbUser, dbPassword, port, "xchan")
	if err != nil {
		c.JSON(200, Response(400, "数据库链接失败"+err.Error(), nil))
		return
	}
	// 链接成功后则进行安装
	// 配置文件保存
	err = h.s.Config.Set("db::host", host)
	if err != nil {
		c.JSON(200, Response(400, "保存配置文件失败"+err.Error(), nil))
		return
	}
	err = h.s.Config.Set("db::name", "xchan")
	if err != nil {
		c.JSON(200, Response(400, "保存配置文件失败"+err.Error(), nil))
		return
	}
	err = h.s.Config.Set("db::user", dbUser)
	if err != nil {
		c.JSON(200, Response(400, "保存配置文件失败"+err.Error(), nil))
		return
	}
	err = h.s.Config.Set("db::password", dbPassword)
	if err != nil {
		c.JSON(200, Response(400, "保存配置文件失败"+err.Error(), nil))
		return
	}
	err = h.s.Config.Set("db::port", port)
	if err != nil {
		c.JSON(200, Response(400, "保存配置文件失败"+err.Error(), nil))
		return
	}

	// 第二步,表同步
	h.s.DB.AutoMigrate(&mod.User{}, &mod.OSS{}, &mod.BaseConfig{})
	// 第三步,判断是哪个存储
	switch position {
	case "local":
		h.s.StoreType = mod.StoreType.Local
		err = h.s.AddLocalConfig(localDomain)
		if err != nil {
			c.JSON(200, Response(400, "保存配置文件失败"+err.Error(), nil))
			return
		}
	case "qiniu":
		h.s.StoreType = mod.StoreType.QiNiu
		err = h.s.AddQiNiuConfig(key, secret, bucket, qiNiuDomain)
		if err != nil {
			c.JSON(200, Response(400, "保存配置文件失败"+err.Error(), nil))
			return
		}
		// 生成管理器
		domain := strings.Trim(qiNiuDomain, "/") + "/"
		h.s.OSS = struct{ Key, Secret, Bucket, Domain string }{Key: key, Secret: secret, Bucket: bucket, Domain: domain}
		mac := qbox.NewMac(h.s.OSS.Key, h.s.OSS.Secret)
		h.s.Manager = storage.NewBucketManager(mac, new(storage.Config))
	}
	// 新增管理员用户
	err = h.s.AddUser(&User{
		UserName: user,
		Password: password,
	})
	if err != nil {
		c.JSON(200, Response(400, "新增管理员失败"+err.Error(), nil))
		return
	}
	// 标记安装完成
	err = h.s.Config.Set("install::have_install", "true")
	if err != nil {
		c.JSON(200, Response(400, "保存配置文件失败"+err.Error(), nil))
		return
	}
	err = h.s.Config.SaveConfigFile("./conf/conf.ini")
	if err != nil {
		c.JSON(200, Response(400, "保存配置文件失败"+err.Error(), nil))
		return
	}

	c.JSON(200, Response(200, "安装成功", nil))
}
