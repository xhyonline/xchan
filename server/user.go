package server

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/xhyonline/xchan/mod"
)

// AddUser 新增一条用户
func (s *Server) AddUser(user *User) error {
	_, exists, err := s.GetUserByUserName(user.UserName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("该用户名已经被注册")
	}
	item := &mod.User{
		Username: user.UserName,
		Password: user.Password,
	}
	err = s.DB.Create(item).Error
	if err != nil {
		return err
	}
	return nil
}

// GetUserByUserName 通过用户名获取一条
func (s *Server) GetUserByUserName(user string) (*mod.User, bool, error) {
	item := new(mod.User)
	err := s.DB.Model(&mod.User{}).Where("username = ?", user).First(item).Error
	// 如果没找到
	if gorm.IsRecordNotFoundError(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return item, true, nil
}

// UpdateUser 通过用户名修改一条
func (s *Server) UpdateUser(user *User) error {
	if user.NewPassword == "" {
		return fmt.Errorf("新密码不能为空")
	}
	if user.NewPassword != user.RePassword {
		return fmt.Errorf("两次密码不一致")
	}
	isTrue, err := s.CheckIDAndPass(user.UserName, user.OldPassword)
	if err != nil {
		return err
	}
	if !isTrue {
		return fmt.Errorf("旧密码错误")
	}
	// 存在则修改
	err = s.DB.Debug().Model(&mod.User{}).Where("username = ? and password = ?", user.UserName, user.OldPassword).
		Update("password", user.NewPassword).Error
	if err != nil {
		return err
	}
	return nil
}
