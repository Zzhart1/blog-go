package service

import (
	"errors"
	"strconv"
	"strings"

	"blog/model"
)

// 可判定的包级错误，controller 据此映射 HTTP 状态码
var (
	ErrInvalidID        = errors.New("无效的 id")
	ErrForbidden        = errors.New("无权限操作该文章")
	ErrCategoryNotFound = errors.New("分类不存在")
)

// parseID 解析正整数 id，非法返回错误
func parseID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		return 0, ErrInvalidID
	}
	return uint(id), nil
}

// checkCategory 校验分类是否存在
func checkCategory(categoryId uint) error {
	var cate model.Category
	if err := model.DB.First(&cate, categoryId).Error; err != nil {
		return ErrCategoryNotFound
	}
	return nil
}

// 创建文章
func CreateArticle(userId uint, title, content string, categoryId uint) (*model.Article, error) {
	if err := checkCategory(categoryId); err != nil {
		return nil, err
	}
	art := model.Article{
		Title:      title,
		Content:    content,
		CategoryID: categoryId,
		UserID:     userId,
	}
	if err := model.DB.Create(&art).Error; err != nil {
		return nil, err
	}
	return &art, nil
}

// 修改文章：只能修改自己写的文章，userId 是登录用户
func UpdateArticle(idStr string, userId uint, title, content string, categoryId uint) error {
	id, err := parseID(idStr)
	if err != nil {
		return err
	}
	var art model.Article
	if err := model.DB.First(&art, id).Error; err != nil {
		return err
	}
	// 权限校验：只能改自己的文章
	if art.UserID != userId {
		return ErrForbidden
	}
	if err := checkCategory(categoryId); err != nil {
		return err
	}
	// Select 指定字段：允许把标题/内容/分类清空（结构体 Updates 默认只更新非零字段）
	return model.DB.Model(&art).
		Select("Title", "Content", "CategoryID").
		Updates(model.Article{Title: title, Content: content, CategoryID: categoryId}).Error
}

// 删除文章：只能删自己的
func DeleteArticle(idStr string, userId uint) error {
	id, err := parseID(idStr)
	if err != nil {
		return err
	}
	var art model.Article
	if err := model.DB.First(&art, id).Error; err != nil {
		return err
	}
	if art.UserID != userId {
		return ErrForbidden
	}
	return model.DB.Delete(&art).Error
}

// 分页查询文章列表，关联作者与分类，支持分类筛选和关键词搜索，按创建时间倒序
func GetArticleList(pageStr, sizeStr, categoryIdStr, keyword string) ([]model.Article, int64, error) {
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	db := model.DB.Model(&model.Article{})
	if cid, err := strconv.ParseUint(categoryIdStr, 10, 64); err == nil && cid > 0 {
		db = db.Where("category_id = ?", cid)
	}
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", like, like)
	}

	// 总数与列表分开查询，避免链式复用同一 Statement
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []model.Article
	err := db.
		Omit("Content"). // 列表不返回正文全文，减小 payload
		Preload("Category").
		Preload("User").
		Order("created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&list).Error
	return list, total, err
}

// 根据 id 获取单篇文章
func GetArticleById(idStr string) (*model.Article, error) {
	id, err := parseID(idStr)
	if err != nil {
		return nil, err
	}
	var art model.Article
	if err := model.DB.Preload("Category").Preload("User").First(&art, id).Error; err != nil {
		return nil, err
	}
	return &art, nil
}
