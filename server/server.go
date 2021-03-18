package server

import (
	"github.com/astaxie/beego/config"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xhyonline/xutil/db"
	"strings"
	"sync"
)

var once sync.Once

var instance *Server

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
		dbConfig, err := c.GetSection("db")
		if err != nil {
			panic(err)
		}

		s.DB = db.NewDataBase(&db.Config{
			Host:          dbConfig["host"],
			Port:          dbConfig["port"],
			User:          dbConfig["user"],
			Password:      dbConfig["password"],
			Name:          dbConfig["name"],
			Lifetime:      3600,
			MaxActiveConn: 30,
			MaxIdleConn:   4,
		})
		oss, err := c.GetSection("oss")
		if err != nil {
			panic(err)
		}
		domain := strings.Trim(oss["domain"], "/") + "/"
		s.OSS = struct{ Key, Secret, Bucket, Domain string }{Key: oss["key"], Secret: oss["secret"], Bucket: oss["bucket"], Domain: domain}
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

// response 错误码 400 错误 200 正确
func response(code int, message string, data interface{}) gin.H {
	return gin.H{
		"code": code,
		"msg":  message,
		"data": data,
	}
}
