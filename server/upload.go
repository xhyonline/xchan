package server

import (
	"context"
	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/storage"
	"mime/multipart"
)

// Upload 文件上传
func (s *Server) Upload(file *multipart.FileHeader) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: s.OSS.Bucket + ":" + s.OSS.Key,
	}
	mac := auth.New(s.OSS.Key, s.OSS.Secret)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": file.Filename,
		},
	}
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	err = formUploader.Put(context.Background(), &ret, upToken, s.OSS.Key, f, file.Size, &putExtra)
	if err != nil {
		return "", err
	}
	return s.OSS.Domain + ret.Key, nil
}
