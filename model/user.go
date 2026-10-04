package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:50;unique" json:"username"`
	Password string `gorm:"size:100" json:"-"` // 序列化时永不输出，防止密码泄露

	// 一对多：一个用户有多篇文章
	Articles []Article `gorm:"foreignKey:UserID" json:"-"`
}
