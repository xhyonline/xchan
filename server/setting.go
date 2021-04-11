package server

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xhyonline/xchan/mod"
	"github.com/xhyonline/xutil/helper"
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
		var flag bool
		var err error
		// 数据库失败尝试 4 次,4次都连不上.....就是失败了
		for i := 0; i < 4; i++ {
			err = s.DB.DB().Ping()
			if err == nil {
				flag = true
				break
			}
			log.Info("数据库链接中断正尝试链接第", i+1, "次")
		}
		if !flag {
			return fmt.Errorf("数据库链接中断 %s", err)
		}
		return nil
	}
	// 如果数据库还没链接,则连
	m, err := s.Config.GetSection("db")
	if err != nil {
		return err
	}
	// 初始化数据库连接
	if err = s.ConnectDB(m["host"], m["user"], m["password"], m["port"], m["name"]); err != nil {
		return err
	}
	// 重载当前配置类型
	return s.reloadStoreType()
}

// reloadStoreType 重载配置类型
func (s *Server) reloadStoreType() error {

	// 先判断有没有设置七牛?
	_, isSet, err := s.ExistsSettingStoryType(mod.StoreType.QiNiu)
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
	_, isSet, err = s.ExistsSettingStoryType(mod.StoreType.Local)
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
func (s *Server) AddQiNiuConfig(data *mod.QiNiuOSSConfig) error {
	if err := data.Validate(); err != nil {
		return err
	}
	s.StoreType = mod.StoreType.QiNiu

	body, err := json.Marshal(data)
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
	s.StoreType = mod.StoreType.Local
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

// UpdateOrAddLocalSetting 修改本地存储设置,或者增加一条配置,还会设置为启动该配置
func (s *Server) UpdateOrAddLocalSetting(domain string) error {
	if !helper.IsURL(domain) {
		return fmt.Errorf("本地绑定的 URL 格式不正确")
	}
	domain = strings.Trim(domain, "/")
	// 如果没有修改最好
	if s.GetLocalDomain() == domain {
		return nil
	}
	// 判断记录是否存在
	item, exists, err := s.ExistsSettingStoryType(mod.StoreType.Local)
	if err != nil {
		return err
	}

	if !exists {
		// 不存在则添加一条,并设置为启动
		return s.AddLocalConfig(domain)
	}
	// 存在则修改原本的配置
	local := new(mod.LocalConfig)
	tmp := "file-save-dir"
	updateDomain := domain + "/" + tmp + "/"
	oldLocalDomain := s.LocalDomain
	s.LocalDomain = updateDomain
	if err := json.Unmarshal([]byte(item.Body), local); err != nil {
		return err
	}

	// 修改域名,默认是不修改存储路径的
	local.Domain = updateDomain

	body, err := json.Marshal(local)
	if err != nil {
		return err
	}
	// 修改本地域名
	item.Body = string(body)
	if err = s.DB.Save(item).Error; err != nil {
		return err
	}

	// 修改原来的图片存储路径,防止以前的图片 404 找不到了
	items, err := s.FindStoreFileByStoreType(mod.StoreType.Local)
	if err != nil {
		return err
	}

	for _, item := range items {
		item.Path = strings.Replace(item.Path, oldLocalDomain, updateDomain, 1)
		if err = s.DB.Save(item).Error; err != nil {
			return err
		}
	}
	// 开启本地存储
	return s.OpenStoreTypeCloseOther(mod.StoreType.Local)
}

// UpdateOrAddQiNiuSetting 修改七牛存储设置,或者增加一条配置,还会设置为启动该配置
func (s *Server) UpdateOrAddQiNiuSetting(data *mod.QiNiuOSSConfig) error {

	if err := data.Validate(); err != nil {
		return err
	}

	domain := strings.Trim(data.Domain, "/") + "/"

	// 判断记录是否存在
	item, exists, err := s.ExistsSettingStoryType(mod.StoreType.QiNiu)
	if err != nil {
		return err
	}
	if !exists {
		// 不存在则添加一条,并设置为启动
		return s.AddQiNiuConfig(data)
	}

	// 如果存在,直接覆盖 body 字段
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	item.Body = string(body)

	if err := s.DB.Save(item).Error; err != nil {
		return err
	}
	items, err := s.FindStoreFileByStoreType(mod.StoreType.QiNiu)
	if err != nil {
		return err
	}
	// 修改存储域名,防止404
	for _, item := range items {
		item.Path = strings.Replace(item.Path, s.OSS.Domain, domain, 1)
		if err := s.DB.Save(item).Error; err != nil {
			return err
		}
	}
	// 开启七牛存储
	return s.OpenStoreTypeCloseOther(mod.StoreType.QiNiu)
}

// OpenStoryType 开启存储类型,其它都关闭
func (s *Server) OpenStoreTypeCloseOther(store mod.StoreTypeEnum) error {
	// 开启该设置
	if err := s.DB.Model(&mod.BaseConfig{}).Where("type = ?", store).Updates(map[string]interface{}{
		"is_open": true,
	}).Error; err != nil {
		return err
	}

	// 不是该存储类型的都关闭
	if err := s.DB.Model(&mod.BaseConfig{}).Where("type != ?", store).Updates(map[string]interface{}{
		"is_open": false,
	}).Error; err != nil {
		return err
	}
	// 重载配置
	return s.reloadStoreType()
}
