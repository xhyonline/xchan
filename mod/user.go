package mod

import "github.com/jinzhu/gorm"

// User 用户表
type User struct {
	gorm.Model
	// 用户名
	Username string
	// 密码
	Password string
}
