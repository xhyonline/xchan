package server

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/xhyonline/xchan/mod"
	"github.com/xhyonline/xutil/helper"
	"mime/multipart"
	"path"
	"strings"
)

// Upload 文件上传
func (s *Server) Upload(file *multipart.FileHeader, user string) (string, error) {
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
		Path: src,
		Size: file.Size,
		User: user,
		Key:  ret.Key,
		Name: file.Filename,
		Hash: hash,
		Ext:  path.Ext(file.Filename),
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
	err = s.Manager.Delete(s.OSS.Bucket, n.Key)
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
