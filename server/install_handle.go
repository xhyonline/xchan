package server

import (
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/xhyonline/xchan/mod"
	"github.com/xhyonline/xutil/helper"
	"strings"
	"sync/atomic"
)

var atomicLock = new(uint32) // 防抖锁

// Install 安装方法
func (h *Handler) Install(c *gin.Context) {
	// 后端防抖,防止手快直接点了两下安装,导致安装了两次,为 1 时代表正在安装
	if atomic.LoadUint32(atomicLock) == 1 {
		// 已经开始了安装程序
		c.JSON(200, Response(400, "正在安装中,请稍后呢~", nil))
		return
	}
	// 安装中
	atomic.StoreUint32(atomicLock, 1)
	// 当然我们不排除安装出错,所以最后还要解锁,让用户配置正确
	defer atomic.StoreUint32(atomicLock, 0)

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

	// 后端防抖

	// 第三步,判断是哪个存储
	switch position {
	case "local":
		err = h.s.AddLocalConfig(localDomain)
		if err != nil {
			c.JSON(200, Response(400, "保存配置文件失败"+err.Error(), nil))
			return
		}
	case "qiniu":
		err = h.s.AddQiNiuConfig(&mod.QiNiuOSSConfig{
			Key:    key,
			Secret: secret,
			Bucket: bucket,
			Domain: qiNiuDomain,
		})
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

func (h *Handler) InstallView(c *gin.Context) {
	log.Infof("用户还没安装")
	c.HTML(200, "install.html", nil)
}
