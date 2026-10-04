package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string    `gorm:"size:50" json:"name"`
	Articles []Article `gorm:"foreignKey:CategoryID" json:"-"` // 关联文章，列表返回时省略
}
