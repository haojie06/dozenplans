package main

import (
	controller "dozenplans/controllers"
	"dozenplans/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func startHttpServer() {
	// 初始化路由
	router := gin.New()
	// 使用中间件，出现错误时写入500
	router.Use(middleware.Metic(), middleware.Cors(), middleware.MyLogger(), gin.Recovery())
	// 分组路由
	// indexGroup := router.Group("/")
	apiGroup := router.Group("/api")
	userGroup := apiGroup.Group("/users")
	testGroup := apiGroup.Group("/test")
	debugGroup := apiGroup.Group("/debug")
	taskGroup := apiGroup.Group("/tasks")
	tagGroup := apiGroup.Group("/tags")
	categoryGroup := apiGroup.Group("/categories")
	progressGroup := apiGroup.Group("/progress")
	// controllGroup := router.Group("/tasks")
	apiGroup.POST("/login", controller.SigninUserHandler)
	{
		debugGroup.POST("/user", controller.AddUser)
		debugGroup.GET("/permission", middleware.Auth(2), controller.NeedPermission)
	}
	{
		userGroup.GET("", controller.GetAllUsersHandler)
		userGroup.GET("/:uid", controller.GetUserHandler)
		userGroup.POST("", controller.RegisterUserHandler)
		userGroup.PUT("", middleware.Auth(0), controller.UpdateUserHandler) // 更新自身信息
	}

	{
		taskGroup.POST("", middleware.Auth(0), controller.AddTaskHandler)
		taskGroup.GET("", middleware.Auth(0), controller.GetTaskHandler)
		taskGroup.GET("/:tid", middleware.Auth(0), controller.GetTaskHandler)
		taskGroup.PUT("/:tid", middleware.Auth(0), controller.UpdateTaskHandler)
		taskGroup.DELETE("/:tid", middleware.Auth(0), controller.DeleteTaskHandler)
	}
	{
		tagGroup.GET("", middleware.Auth(0), controller.GetAllTagHandler)
		tagGroup.GET("/:id", middleware.Auth(0), controller.GetTasksByTag)
	}
	{
		categoryGroup.GET("/:id", middleware.Auth(0), controller.GetTasksByCategory)
		categoryGroup.GET("", middleware.Auth(0), controller.GetAllCategoriesByUid)
		// categoryGroup.POST("", middleware.Auth(0), controller.CreateCategory)
	}
	{
		progressGroup.GET("", middleware.Auth(0), controller.GetAllProgress)
	}
	{
		testGroup.GET("/*visit", controller.DemoHandler)
	}
	// 监听并在 0.0.0.0:8080 上启动服务
	log.Println("Http server is running")
	router.Run("127.0.0.1:8080")
}
