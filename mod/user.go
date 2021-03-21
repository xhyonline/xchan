package mod

import (
	"github.com/xhyonline/xutil/model"
)

// User 用户表
type User struct {
	model.Addon
	// 用户名
	Username string
	// 密码
	Password string
}

// OSS 对象存储
type OSS struct {
	model.Addon
	// 存储路径
	Path string `gorm:"index"`
	// OSS key
	Key string
	// 文件大小
	Size int64
	// 操作者
	User string
}
