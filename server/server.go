package server

import "C"
import (
	"github.com/astaxie/beego/config"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/xhyonline/xchan/mod"
	"github.com/xhyonline/xutil/xlog"
	"sync"
)

var once sync.Once

var instance *Server

var log = xlog.Get(true)

// Server
type Server struct {
	// 数据库实例
	DB *gorm.DB
	// 配置信息
	Config config.Configer
	// OSS 对象存储
	OSS struct {
		Key, Secret, Bucket, Domain string
	}
	// 七牛云对象存储管理者,它负责文件的删除,但不负责文件的上传
	Manager *storage.BucketManager
	// 如果是本地存储,这里将会存存储路径
	PathDir, LocalDomain string
	// 存储类型
	StoreType mod.StoreTypeEnum
}

// GetService 获取标准服务
func GetService() *Server {
	once.Do(func() {
		s := new(Server)
		c, err := config.NewConfig("ini", "conf/conf.ini")
		if err != nil {
			panic(err)
		}
		s.Config = c
		instance = s
	})
	return instance
}

// Handler
type Handler struct {
	s *Server
}

// NewHandler
func NewHandler(s *Server) *Handler {
	return &Handler{s: s}
}

// Response 错误码 400 错误 200 正确
func Response(code int, message string, data interface{}) gin.H {
	return gin.H{
		"code": code,
		"msg":  message,
		"data": data,
	}
}
