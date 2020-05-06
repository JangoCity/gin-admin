package api

import (
	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/middleware"
	"github.com/LyricTian/gin-admin/internal/app/routers/api/ctl"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/gin-gonic/gin"
)

// RegisterRouter 注册/api路由
func RegisterRouter(app *gin.Engine, b *bll.Common, a auth.Auther) {
	g := app.Group("/api")

	// 用户身份授权
	g.Use(middleware.UserAuthMiddleware(
		a,
		middleware.AllowMethodAndPathPrefixSkipper(
			middleware.JoinRouter("GET", "/api/v1/login"),
			middleware.JoinRouter("POST", "/api/v1/login"),
		),
	))

	// 请求频率限制中间件
	g.Use(middleware.RateLimiterMiddleware())

	demoCtl := ctl.NewDemo(b)

	v1 := g.Group("/v1")
	{
		// 注册/api/v1/demos
		v1.GET("/demos", demoCtl.Query)
		v1.GET("/demos/:id", demoCtl.Get)
		v1.POST("/demos", demoCtl.Create)
		v1.PUT("/demos/:id", demoCtl.Update)
		v1.DELETE("/demos/:id", demoCtl.Delete)
		v1.PATCH("/demos/:id/enable", demoCtl.Enable)
		v1.PATCH("/demos/:id/disable", demoCtl.Disable)
	}
}
