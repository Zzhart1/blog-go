package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content   string `gorm:"type:text" json:"content"`
	ArticleID uint   `json:"article_id"`
	UserID    uint   `json:"user_id"`

	// 关联：评论者（序列化时不含密码）
	User User `gorm:"foreignKey:UserID" json:"user"`
}
