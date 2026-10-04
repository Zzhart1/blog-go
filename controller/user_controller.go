package controller

import (
	"errors"
	"net/http"

	"blog/middleware"
	"blog/service"

	"github.com/gin-gonic/gin"
)

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误：用户名3-20位，密码6-72位"})
		return
	}
	// 调用 service 层，controller 不再直接操作 DB
	_, err := service.RegisterUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"msg": "用户名已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "注册失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
}

func Login(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	user, err := service.LoginUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrBadLogin) {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "账号密码错误"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "登录失败"})
		return
	}
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "登录失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "msg": "登陆成功"})
}
