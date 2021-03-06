package mod

import (
	"fmt"
	"github.com/xhyonline/xutil/helper"
	"github.com/xhyonline/xutil/model"
)

// StoreType 存储类型
type StoreTypeEnum int

var StoreType = struct {
	Local StoreTypeEnum
	QiNiu StoreTypeEnum
}{
	// 本地
	Local: 1,
	// 七牛
	QiNiu: 2,
}

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
	// 存储路径是一个 URL
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
	// 存储类型
	StoreType StoreTypeEnum `gorm:"tinyint(1)"`
	// 本地文件存放地址,注意七牛云是没有的
	LocalFilePath string
}

// BaseConfig 基础配置
type BaseConfig struct {
	model.Addon
	// 1.本地存储、2、七牛存储
	Type StoreTypeEnum `gorm:"tinyint(1)"`
	// 具体配置内容
	Body string `gorm:"varchar(600)"`
	// 是否开启
	IsOpen bool
}

// QiNiuOSSConfig 七牛对象存储配置 ,不是数据表
type QiNiuOSSConfig struct {
	// 七牛云 KEY
	Key string `json:"key"`
	// 七牛云 Secret
	Secret string `json:"secret"`
	// 七牛云存储同
	Bucket string `json:"bucket"`
	// 绑定的七牛云域名
	Domain string `json:"domain"`
}

func (data *QiNiuOSSConfig) Validate() error {
	if data.Key == "" || data.Domain == "" || data.Secret == "" || data.Bucket == "" {
		return fmt.Errorf("请填写完整配置喔")
	}
	if !helper.IsURL(data.Domain) {
		return fmt.Errorf("七牛云绑定的 URL 格式不正确")
	}
	return nil
}

// LocalConfig 本地存储配置,,不是数据表
type LocalConfig struct {
	// 本地存储域名
	Domain string `json:"domain"`
	// 存储的绝对路径目录
	Path string `json:"path"`
}
