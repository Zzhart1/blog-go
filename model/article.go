package model

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title      string `gorm:"size:100" json:"title"`
	Content    string `gorm:"type:text" json:"content"`
	CategoryID uint   `json:"category_id"` // 分类 id 外键
	UserID     uint   `json:"user_id"`     // 作者 id

	// 关联（BelongsTo）
	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
	User     User     `gorm:"foreignKey:UserID" json:"author"`
}
