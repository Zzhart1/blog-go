package service

import (
	"errors"
	"strings"

	"blog/model"
)

// ErrCommentEmpty 评论内容为空
var ErrCommentEmpty = errors.New("评论内容不能为空")

// CreateComment 发表评论：校验文章存在、内容非空
func CreateComment(userId uint, articleId uint, content string) (*model.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, ErrCommentEmpty
	}

	var art model.Article
	if err := model.DB.First(&art, articleId).Error; err != nil {
		return nil, err
	}

	comment := model.Comment{
		Content:   content,
		ArticleID: articleId,
		UserID:    userId,
	}
	if err := model.DB.Create(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetCommentsByArticle 按文章查评论，按创建时间正序，带出评论者信息
func GetCommentsByArticle(articleIdStr string) ([]model.Comment, error) {
	id, err := parseID(articleIdStr)
	if err != nil {
		return nil, err
	}
	var list []model.Comment
	err = model.DB.
		Where("article_id = ?", id).
		Preload("User").
		Order("created_at ASC").
		Find(&list).Error
	return list, err
}

// DeleteComment 删除评论：只能删自己的
func DeleteComment(idStr string, userId uint) error {
	id, err := parseID(idStr)
	if err != nil {
		return err
	}
	var comment model.Comment
	if err := model.DB.First(&comment, id).Error; err != nil {
		return err
	}
	if comment.UserID != userId {
		return ErrForbidden
	}
	return model.DB.Delete(&comment).Error
}
