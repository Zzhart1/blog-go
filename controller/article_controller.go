package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"blog/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 创建文章请求结构体
type CreateArticleReq struct {
	Title      string `json:"title" binding:"required,max=100"`
	Content    string `json:"content" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
}

type UpdateArticleReq struct {
	Title      string `json:"title" binding:"required,max=100"`
	Content    string `json:"content" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
}

// 创建文章
func CreateArticle(c *gin.Context) {
	var req CreateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	userId := c.GetUint("userId")
	art, err := service.CreateArticle(userId, req.Title, req.Content, req.CategoryID)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "分类不存在"})
			return
		}
		slog.Error("创建文章失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "创建成功", "data": art})
}

// 修改文章
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var req UpdateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	userId := c.GetUint("userId")
	err := service.UpdateArticle(id, userId, req.Title, req.Content, req.CategoryID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的文章 id"})
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"msg": "无权修改该文章"})
		case errors.Is(err, service.ErrCategoryNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"msg": "分类不存在"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"msg": "文章不存在"})
		default:
			slog.Error("修改文章失败", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "修改失败"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "修改成功"})
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetUint("userId")
	err := service.DeleteArticle(id, userId)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的文章 id"})
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"msg": "无权删除该文章"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"msg": "文章不存在"})
		default:
			slog.Error("删除文章失败", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

// 获取文章列表（分页 + 分类筛选 + 关键词搜索）
func GetArticleList(c *gin.Context) {
	page := c.Query("page")
	size := c.Query("size")
	categoryId := c.Query("category_id")
	keyword := c.Query("keyword")
	list, total, err := service.GetArticleList(page, size, categoryId, keyword)
	if err != nil {
		slog.Error("查询文章列表失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  list,
		"total": total,
	})
}

// 获取单篇文章详情
func GetArticleDetail(c *gin.Context) {
	id := c.Param("id")
	art, err := service.GetArticleById(id)
	if err != nil {
		if errors.Is(err, service.ErrInvalidID) {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的文章 id"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"msg": "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": art})
}
