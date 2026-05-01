package routes

import (
	api2 "htmlhub/api"
	middleware2 "htmlhub/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 跨域中间件（放在最前面）
	r.Use(middleware2.Cors())
	// 注册全局异常处理中间件
	r.Use(middleware2.ErrorHandler())

	userApi := api2.UserApi{}
	// 用户路由
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/login", userApi.Login)
		userGroup.POST("/register", userApi.Register)
		userGroup.GET("/test", userApi.Test)
	}

	exampleApi := api2.ExampleApi{}
	htmlRecordApi := api2.HTMLRecordApi{}
	htmlRecordDataApi := api2.HTMLRecordDataApi{}
	r.GET("/", htmlRecordApi.PublicHTML)
	// 用户子域上误打开 /home、/index.html 时仍返回注入后的 HTML（与 / 一致，避免落到 Vue 的 history 路由）
	r.GET("/home", htmlRecordApi.PublicHTML)
	r.GET("/index.html", htmlRecordApi.PublicHTML)
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware2.JWTInterceptor()) // 应用JWT拦截器
	{
		// 消费记录路由（需要认证）
		exampleGroup := apiGroup.Group("/example")
		{
			exampleGroup.POST("/test", exampleApi.Test) // 添加消费记录
		}
		// 消费拓展路由（需要认证）
		htmlGroup := apiGroup.Group("/html")
		{
			htmlGroup.POST("/upload", htmlRecordApi.Upload)
			htmlGroup.GET("/my", htmlRecordApi.MyList)
			htmlGroup.DELETE("/:id", htmlRecordApi.Delete)
			htmlGroup.PUT("/:id/visibility", htmlRecordApi.UpdateVisibility)
			htmlGroup.POST("/data/save", htmlRecordDataApi.Save)
			htmlGroup.GET("/data/load", htmlRecordDataApi.Load)
		}

	}

	return r
}
