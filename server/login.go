package server

import (
	"github.com/xhyonline/xchan/mod"
)

// CheckLogin 登录校验
func (s *Server) CheckLogin(id, pass string) (bool, error) {
	var count int
	err := s.DB.Debug().Model(&mod.User{}).Where("username = ? and password = ? ", id, pass).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
