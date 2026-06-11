package routes

import (
	"htmlhub/api"
	"htmlhub/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 跨域中间件（放在最前面）
	r.Use(middleware.Cors())
	// 注册全局异常处理中间件
	r.Use(middleware.ErrorHandler())

	userApi := api.UserApi{}
	// 用户路由
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/login", middleware.LoginIPRateLimit(), userApi.Login)
		userGroup.POST("/register", middleware.LoginIPRateLimit(), userApi.Register)
		userGroup.GET("/test", userApi.Test)
	}

	exampleApi := api.ExampleApi{}
	htmlRecordApi := api.HTMLRecordApi{}
	htmlRecordDataApi := api.HTMLRecordDataApi{}
	htmlRecordLikeApi := api.HTMLRecordLikeApi{}
	//r.GET("/", htmlRecordApi.PublicHTML)
	// 用户子域上误打开 /home、/index.html 时仍返回注入后的 HTML（与 / 一致，避免落到 Vue 的 history 路由）
	//r.GET("/home", htmlRecordApi.PublicHTML)
	//r.GET("/index.html", htmlRecordApi.PublicHTML)

	r.GET("/", middleware.HighRiskWriteRateLimit(), htmlRecordApi.PublicHTML)
	r.GET("/home", middleware.HighRiskWriteRateLimit(), htmlRecordApi.PublicHTML)
	r.GET("/index.html", middleware.HighRiskWriteRateLimit(), htmlRecordApi.PublicHTML)
	r.GET("/api/html/public", htmlRecordApi.PublicList)
	r.GET("/api/html/data/publish", middleware.HighRiskWriteRateLimit(), htmlRecordDataApi.PublishLoad)
	r.GET("/api/html/:id/likes/count", htmlRecordLikeApi.Count)

	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.JWTInterceptor()) // 应用JWT拦截器
	{
		// 路由（需要认证）
		exampleGroup := apiGroup.Group("/example")
		{
			exampleGroup.POST("/test", exampleApi.Test)
		}
		// 当前用户自助路由（个人中心）
		userAuthGroup := apiGroup.Group("/user")
		{
			userAuthGroup.GET("/profile", userApi.Profile)
			userAuthGroup.PUT("/profile", userApi.UpdateProfile)
		}
		// 路由（需要认证）
		htmlGroup := apiGroup.Group("/html")
		htmlGroup.Use(middleware.HighRiskWriteRateLimit())
		{
			htmlGroup.POST("/upload", htmlRecordApi.Upload)
			htmlGroup.GET("/my", htmlRecordApi.MyList)
			htmlGroup.GET("/liked", htmlRecordLikeApi.MyLikes)
			htmlGroup.POST("/data/save", htmlRecordDataApi.Save)
			htmlGroup.GET("/data/load", htmlRecordDataApi.Load)
			htmlGroup.GET("/data/my", htmlRecordDataApi.MyDataList)
			htmlGroup.DELETE("/data/my", htmlRecordDataApi.ClearMyData)
			htmlGroup.GET("/data/my/:id/export", htmlRecordDataApi.ExportMyData)
			htmlGroup.DELETE("/data/my/:id", htmlRecordDataApi.DeleteMyData)
			htmlGroup.DELETE("/:id", htmlRecordApi.Delete)
			htmlGroup.PUT("/:id/description", htmlRecordApi.UpdateDescription)
			htmlGroup.PUT("/:id/content", htmlRecordApi.UpdateHTMLContent)
			htmlGroup.PUT("/:id/visibility", htmlRecordApi.UpdateVisibility)
			htmlGroup.PUT("/:id/publish-mode", htmlRecordApi.UpdatePublishMode)
			htmlGroup.POST("/:id/like", htmlRecordLikeApi.Like)
			htmlGroup.DELETE("/:id/like", htmlRecordLikeApi.Unlike)
		}
		adminGroup := apiGroup.Group("/admin")
		adminGroup.Use(middleware.AdminInterceptor())
		adminGroup.Use(middleware.HighRiskWriteRateLimit())
		{
			adminGroup.GET("/users", userApi.AdminList)
			adminGroup.PUT("/users/:id", userApi.AdminUpdateUser)
			adminGroup.GET("/html", htmlRecordApi.AdminList)
			adminGroup.GET("/html/:id", htmlRecordApi.AdminDetail)
			adminGroup.PUT("/html/:id/approval", htmlRecordApi.AdminUpdateApprovalStatus)
			adminGroup.PUT("/html/:id/subdomain", htmlRecordApi.AdminUpdateSubdomain)
			adminGroup.DELETE("/html/:id", htmlRecordApi.AdminDelete)
			adminGroup.GET("/html-data", htmlRecordDataApi.AdminList)
		}

	}

	return r
}
