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
	// 文件名
	Name string `json:"name"`
	// 存储路径
	Path string `gorm:"index" json:"path"`
	// OSS key
	Key string `json:"key"`
	// 文件大小
	Size int64 `json:"size" `
	// 操作者
	User string `json:"user"`
	// hash 校验和
	Hash string `json:"hash"`
	// 后缀名
	Ext string `json:"ext"`
	// 辅助字段单位不在数据库中存储
	Unit string `json:"unit" gorm:"-" `
	// 辅助字段转换不在数据库中存储
	SizeFormat string `json:"size_format" gorm:"-"`
	// 辅助字段转换不在数据库中存储
	TimeFormat string `json:"time" gorm:"-"`
}
