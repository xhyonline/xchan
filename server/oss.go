package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/xhyonline/xchan/mod"
	"github.com/xhyonline/xutil/helper"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// UploadQiNiu 文件上传
func (s *Server) UploadQiNiu(file *multipart.FileHeader, user string) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: s.OSS.Bucket,
	}
	mac := auth.New(s.OSS.Key, s.OSS.Secret)
	upToken := putPolicy.UploadToken(mac)

	formUploader := storage.NewFormUploader(new(storage.Config))
	ret := storage.PutRet{}

	putExtra := &storage.PutExtra{
		Params: map[string]string{
			"x:name": file.Filename,
		},
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, file.Size)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}
	// 计算一次 hash
	tmp := string(buffer)
	hash := helper.Md5(tmp)

	// 本地查询是否存在
	item, exists, err := s.Exists(hash)
	if err != nil {
		return "", err
	}
	if exists {
		return item.Path, nil
	}
	// 不存在则上传
	r := strings.NewReader(tmp)
	err = formUploader.Put(context.Background(), &ret, upToken, hash, r, file.Size, putExtra)
	if err != nil {
		return "", err
	}

	src := s.OSS.Domain + ret.Key

	// 入库
	err = s.DB.Create(&mod.OSS{
		Path:      src,
		Size:      file.Size,
		User:      user,
		Key:       ret.Key,
		Name:      file.Filename,
		Hash:      hash,
		Ext:       path.Ext(file.Filename),
		StoreType: mod.StoreType.QiNiu,
	}).Error

	if err != nil {
		return "", err
	}
	return src, nil
}

// UploadLocal 上传文件到本地
func (s *Server) UploadLocal(file *multipart.FileHeader, user string) (string, error) {
	// 查看目录是否存在
	exists, err := helper.PathExists(s.PathDir)
	log.Infof("存储路径:%s", s.PathDir)
	if err != nil {
		log.Errorf("判断存储路径是否存在出错 %s:", err)
		return "", err
	}
	if !exists {
		err = os.MkdirAll(s.PathDir, 777)
		if err != nil {
			log.Errorf("创建存储路径出错 %s:", err)
			return "", err
		}
	}
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, file.Size)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}
	// 计算一次 hash
	strBuffer := string(buffer)
	hash := helper.Md5(strBuffer)
	// 本地查询是否存在
	item, exists, err := s.Exists(hash)
	if err != nil {
		return "", err
	}
	if exists {
		return item.Path, nil
	}
	// 后缀
	ext := path.Ext(file.Filename)
	filePath := s.PathDir + `/` + hash + ext
	if strings.Contains(filePath, `\`) {
		filePath = s.PathDir + `\` + hash + ext
	}
	// 不存在则上传
	newFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	reader := strings.NewReader(strBuffer)
	if _, err = io.Copy(newFile, reader); err != nil {
		return "", err
	}
	src := s.LocalDomain + hash + ext
	// 创建文件之后开始入库
	// 入库
	err = s.DB.Create(&mod.OSS{
		Path:          src,
		Size:          file.Size,
		User:          user,
		Key:           hash,
		Name:          file.Filename,
		Hash:          hash,
		Ext:           ext,
		StoreType:     mod.StoreType.Local,
		LocalFilePath: filePath,
	}).Error

	if err != nil {
		return "", err
	}
	return src, nil
}

// Exists 是否存在
func (s *Server) Exists(hash string) (*mod.OSS, bool, error) {
	first := new(mod.OSS)
	err := s.DB.Where("hash =  ?", hash).First(first).Error
	// 如果记录找到了
	if gorm.IsRecordNotFoundError(err) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return first, true, nil
}

// Remove 删除文件
func (s *Server) Remove(id string) error {
	n, err := s.GetByID(id)
	if err != nil {
		return err
	}
	switch n.StoreType {
	case mod.StoreType.Local: // 本地存储
		err = os.Remove(n.LocalFilePath)
	case mod.StoreType.QiNiu: // 七牛存储
		err = s.Manager.Delete(s.OSS.Bucket, n.Key)
	}
	if err != nil {
		return err
	}

	// 删除数据库的数据
	err = s.DB.Debug().Where("id = ?", id).Delete(&mod.OSS{}).Error
	if err != nil {
		return err
	}
	return nil
}

// RemoveList 删除文件
func (s *Server) RemoveList(ids []string) error {
	for _, v := range ids {
		if err := s.Remove(v); err != nil {
			return err
		}
	}
	return nil
}

// GetByPath 通过路径获取一条记录
func (s *Server) GetByPath(path string) (*mod.OSS, error) {
	n := new(mod.OSS)
	err := s.DB.Model(&mod.OSS{}).Where("path = ?", path).First(n).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("没有这条记录")
	}
	if err != nil {
		return nil, err
	}
	return n, nil
}

// GetByID 通过路径获取一条记录
func (s *Server) GetByID(id string) (*mod.OSS, error) {
	n := new(mod.OSS)
	err := s.DB.Model(&mod.OSS{}).Where("id = ?", id).First(n).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("没有这条记录")
	}
	if err != nil {
		return nil, err
	}
	return n, nil
}

// GetListByPaging 分页
func (s *Server) GetListByPaging(page, limit int) ([]*mod.OSS, int, error) {
	offset := (page - 1) * limit
	resp := make([]*mod.OSS, 0)
	err := s.DB.Where(&mod.OSS{}).Offset(offset).Limit(limit).Order("created_at desc").Find(&resp).Error
	if err != nil {
		return nil, 0, err
	}
	for _, v := range resp {
		v.SizeFormat, v.Unit = FormatFileSizeAndUnit(v.Size)
		v.TimeFormat = helper.TimeStampToDate(int(v.CreatedAt.Unix()))
	}
	total, err := s.GetTotal()
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}

// GetTotal 获取总数
func (s *Server) GetTotal() (int, error) {
	var count int
	err := s.DB.Model(&mod.OSS{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SetOSSManager 设置对象存储管理器
func (s *Server) SetOSSManager() error {
	base := new(mod.BaseConfig)
	// 获取一个
	err := s.DB.Model(&mod.BaseConfig{}).Where("type = ?", mod.StoreType.QiNiu).First(base).Error
	if err != nil {
		return err
	}
	oss := new(mod.QiNiuOSSConfig)
	err = json.Unmarshal([]byte(base.Body), oss)
	if err != nil {
		return err
	}
	// 生成管理器
	domain := strings.Trim(oss.Domain, "/") + "/"
	s.OSS = struct{ Key, Secret, Bucket, Domain string }{Key: oss.Key, Secret: oss.Secret, Bucket: oss.Bucket, Domain: domain}
	mac := qbox.NewMac(s.OSS.Key, s.OSS.Secret)
	s.Manager = storage.NewBucketManager(mac, new(storage.Config))
	return nil
}

// SetLocalStorePath 设置本地存储路径
func (s *Server) SetLocalStorePath() error {
	base := new(mod.BaseConfig)
	// 获取一个配置
	err := s.DB.Model(&mod.BaseConfig{}).Where("type = ?", mod.StoreType.Local).First(base).Error
	if err != nil {
		return err
	}
	local := new(mod.LocalConfig)
	err = json.Unmarshal([]byte(base.Body), local)
	if err != nil {
		return err
	}
	s.PathDir = local.Path
	s.LocalDomain = local.Domain
	return nil
}

// GetStoryType 获取存储类型
func (s *Server) GetStoryType() (mod.StoreTypeEnum, error) {
	base := new(mod.BaseConfig)
	err := s.DB.Model(&mod.BaseConfig{}).Where("is_open = ?", true).First(base).Error
	if err != nil {
		return 0, err
	}
	return base.Type, nil
}

// ExistsSettingStoryType 判断用户是否有设置过某种存储类型
func (s *Server) ExistsSettingStoryType(storeType mod.StoreTypeEnum) (*mod.BaseConfig, bool, error) {
	base := new(mod.BaseConfig)

	err := s.DB.Model(&mod.BaseConfig{}).Where("type = ?", storeType).First(base).Error
	// 没找到
	if gorm.IsRecordNotFoundError(err) {
		return nil, false, nil
	}
	return base, err == nil, err
}

// FindStoreFileByStoreType 获取文件
func (s *Server) FindStoreFileByStoreType(store mod.StoreTypeEnum) ([]*mod.OSS, error) {
	resp := make([]*mod.OSS, 0)
	return resp, s.DB.Model(&mod.OSS{}).Where("store_type = ?", store).Find(&resp).Error
}
