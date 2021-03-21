package server

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/xhyonline/xchan/mod"
	"mime/multipart"
)

// Upload 文件上传
func (s *Server) Upload(file *multipart.FileHeader, user string) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: s.OSS.Bucket + ":" + s.OSS.Key,
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
	err = formUploader.Put(context.Background(), &ret, upToken, s.OSS.Key, f, file.Size, putExtra)
	if err != nil {
		return "", err
	}
	path := s.OSS.Domain + ret.Key

	exists, err := s.Exists(path)
	if err != nil {
		return "", err
	}

	if exists {
		return path, nil
	}
	// 不存在则入库
	err = s.DB.Create(&mod.OSS{
		Path: path,
		Size: file.Size,
		User: user,
		Key:  ret.Key,
	}).Error

	if err != nil {
		return "", err
	}

	return path, nil
}

// Exists 是否存在
func (s *Server) Exists(path string) (bool, error) {
	first := new(mod.OSS)
	err := s.DB.Where("path =  ?", path).First(first).Error
	// 如果记录找到了
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Remove 删除文件
func (s *Server) Remove(path string) error {
	n, err := s.GetByPath(path)
	if err != nil {
		return err
	}
	err = s.Manager.Delete(s.OSS.Bucket, n.Key)
	if err != nil {
		return err
	}
	// 删除数据库的数据
	err = s.DB.Debug().Where("path = ?", path).Delete(&mod.OSS{}).Error
	if err != nil {
		return err
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
