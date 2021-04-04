package server

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xhyonline/xchan/mod"
	"time"
)

// ConnectDB 链接数据库
func (s *Server) ConnectDB(host, user, pass, port, dbname string) error {
	db, err := gorm.Open("mysql", user+":"+pass+
		"@tcp("+host+":"+port+")/"+dbname+
		"?charset=utf8mb4&parseTime=True&loc=Local&timeout=90s")
	if err != nil {
		return err
	}
	err = db.DB().Ping()
	if err != nil {
		return err
	}
	db.DB().SetConnMaxLifetime(3600 * time.Second)
	db.DB().SetMaxOpenConns(30)
	db.DB().SetMaxIdleConns(20)
	s.DB = db
	return nil
}

// AddQiNiuConfig 新增一条七牛存储配置信息
func (s *Server) AddQiNiuConfig(key, secret, bucket, domain string) error {
	t := &mod.OSSConfig{
		Key:    key,
		Secret: secret,
		Bucket: bucket,
		Domain: domain,
	}
	body, err := json.Marshal(t)
	if err != nil {
		return err
	}
	// 1.本地存储、2、七牛存储
	err = s.DB.Create(&mod.BaseConfig{
		Type:   mod.StoreType.QiNiu,
		Body:   string(body),
		IsOpen: true,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// AddLocalConfig 新增本地存储配置
func (s *Server) AddLocalConfig() error {
	path, err := GetCurrentPath()
	if err != nil {
		return err
	}
	var tmp = "file-save-dir"
	// 路径位置
	conf := &mod.LocalConfig{Path: path + tmp}

	body, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	// 创建文件存储配置路径
	err = s.DB.Model(&mod.BaseConfig{}).Create(&mod.BaseConfig{
		Type:   mod.StoreType.Local,
		Body:   string(body),
		IsOpen: true,
	}).Error
	return err
}
