package server

import (
	"github.com/jinzhu/gorm"
	"github.com/xhyonline/xchan/mod"
)

// CheckIDAndPass 登录校验
func (s *Server) CheckIDAndPass(id, pass string) (bool, error) {
	item := new(mod.User)
	err := s.DB.Model(&mod.User{}).Where("username = ? and password = ? ", id, pass).First(item).Error
	// 没有这个用户
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
