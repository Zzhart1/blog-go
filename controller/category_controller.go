package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"blog/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateCategoryReq struct {
	Name string `json:"name" binding:"required,max=50"`
}

// 创建分类
func CreateCategory(c *gin.Context) {
	var req CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	cate, err := service.CreateCategory(req.Name)
	if err != nil {
		slog.Error("创建分类失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建分类失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "创建成功", "data": cate})
}

// 获取分类列表
func GetCategoryList(c *gin.Context) {
	list, err := service.GetCategoryList()
	if err != nil {
		slog.Error("查询分类失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询分类失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// 修改分类
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	cate, err := service.UpdateCategory(id, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的分类 id"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"msg": "分类不存在"})
		default:
			slog.Error("修改分类失败", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "修改分类失败"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "修改成功", "data": cate})
}

// 删除分类
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	err := service.DeleteCategory(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的分类 id"})
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"msg": "分类不存在"})
		default:
			slog.Error("删除分类失败", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
