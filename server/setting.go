package server

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xhyonline/xchan/mod"
	"strings"
	"time"
)

// ConnectDB 链接数据库
func (s *Server) ConnectDB(host, user, pass, port, dbname string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=90s", user, pass, host, port, dbname)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Errorf("链接数据库失败 %s dsn:%s", err, dsn)
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

// CheckInstalled  安装后的初始化检查
func (s *Server) CheckInstalled() error {
	// 如果数据库已经链接了,并且存储类型也设置了,我还是以防万一,再 ping 一次
	if s.DB != nil && s.StoreType != 0 {
		return s.DB.DB().Ping()
	}
	// 尝试链接
	m, err := s.Config.GetSection("db")
	if err != nil {
		return err
	}
	// 初始化数据库连接
	if err = s.ConnectDB(m["host"], m["user"], m["password"], m["port"], m["name"]); err != nil {
		return err
	}
	// 设置存储类型
	if err = s.setStoreType(); err != nil {
		return err
	}
	return nil
}

// setStoreType 设置存储类型
func (s *Server) setStoreType() error {
	// 如果已经设置存储类型了
	if s.StoreType != 0 {
		return nil
	}
	// 反之先设置存储类型
	// 先判断有没有设置七牛?
	isSet, err := s.ExistsSettingStoryType(mod.StoreType.QiNiu)
	if err != nil {
		return err
	}
	// 如果设置了,初始化七牛管理器
	if isSet {
		err = s.SetOSSManager()
		if err != nil {
			return err
		}
	}
	// 在判断是否设置了本地存储
	isSet, err = s.ExistsSettingStoryType(mod.StoreType.Local)
	if err != nil {
		return err
	}
	// 如果设置了,初始化本地路径
	if isSet {
		err = s.SetLocalStorePath()
		if err != nil {
			return err
		}
	}
	log.Info("路径:", s.PathDir)
	// 最后获取存储类型,并设置
	s.StoreType, err = s.GetStoryType()
	if err != nil {
		return err
	}
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
func (s *Server) AddLocalConfig(domain string) error {
	path, err := GetCurrentPath()
	if err != nil {
		return err
	}
	var tmp = "file-save-dir"
	domain = strings.Trim(domain, "/") + "/" + tmp + "/"
	s.LocalDomain = domain
	// 路径位置
	conf := &mod.LocalConfig{Path: path + tmp, Domain: domain}
	s.PathDir = path + tmp
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
