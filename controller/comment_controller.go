package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"blog/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateCommentReq struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// 发表评论（需登录）
func CreateComment(c *gin.Context) {
	var req CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	userId := c.GetUint("userId")
	comment, err := service.CreateComment(userId, req.ArticleID, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCommentEmpty):
			c.JSON(http.StatusBadRequest, gin.H{"msg": "评论内容不能为空"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"msg": "文章不存在"})
		default:
			slog.Error("发表评论失败", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "评论失败"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "评论成功", "data": comment})
}

// 获取某篇文章的评论列表（公开）
func GetCommentList(c *gin.Context) {
	id := c.Param("id")
	list, err := service.GetCommentsByArticle(id)
	if err != nil {
		if errors.Is(err, service.ErrInvalidID) {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的文章 id"})
			return
		}
		slog.Error("查询评论失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询评论失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// 删除评论（需登录，仅作者）
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetUint("userId")
	err := service.DeleteComment(id, userId)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的评论 id"})
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"msg": "无权删除该评论"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"msg": "评论不存在"})
		default:
			slog.Error("删除评论失败", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除评论失败"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
