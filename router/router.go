package router

import (
	"blog/controller"
	"blog/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 不信任代理转发头，防止客户端伪造 IP（X-Forwarded-For 等）
	_ = r.SetTrustedProxies(nil)

	// 跨域中间件（前端分离部署时使用，同源部署不受影响）
	r.Use(middleware.CORS())

	// 公开接口
	public := r.Group("/api")
	{
		public.POST("/register", controller.Register)
		public.POST("/login", middleware.LoginRateLimit(), controller.Login)
		public.GET("/article/list", controller.GetArticleList)
		public.GET("/article/:id", controller.GetArticleDetail)
		public.GET("/article/:id/comments", controller.GetCommentList)
		public.GET("/category/list", controller.GetCategoryList)
	}

	// 需要登录鉴权接口
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/article", controller.CreateArticle)
		auth.PUT("/article/:id", controller.UpdateArticle)
		auth.DELETE("/article/:id", controller.DeleteArticle)
		auth.POST("/comment", controller.CreateComment)
		auth.DELETE("/comment/:id", controller.DeleteComment)
		auth.POST("/category", controller.CreateCategory)
		auth.PUT("/category/:id", controller.UpdateCategory)
		auth.DELETE("/category/:id", controller.DeleteCategory)
	}
	return r
}
