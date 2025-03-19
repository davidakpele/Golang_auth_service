package routers

import (
    "api-service/controllers"
    "api-service/middleware"
    "github.com/gin-gonic/gin"
)

// RegisterRoutes initializes the routes for the application
func RegisterRoutes(router *gin.Engine, 
	authController *controllers.AuthController, 
	userController *controllers.UserController,
	bookmarkController *controllers.BookmarkController,
	commentController *controllers.CommentController,
	reportController *controllers.ReportController,
	resourceController *controllers.ResourceController,
	adminController *controllers.AdminController) {

 	public := router.Group("/auth")
	{
		public.POST("/login", authController.Login)
		public.POST("/register", authController.Register)
		public.GET("/logout", authController.Logout)
		public.POST("/verify-account", authController.VerifyAccount)
		public.POST("/resend-otp", authController.ResendOTP)
	}

	// Private routes
	private := router.Group("/api")
	private.Use(middleware.AuthenticationMiddleware())
	{
		private.GET("/user/:id", userController.GetUserByID) 
		private.DELETE("/user/:id", userController.Delete) 
	}
	commentPrivate := router.Group("/comment")
	commentPrivate.Use(middleware.AuthenticationMiddleware())
	{
		commentPrivate.POST("/create", commentController.CreateComment)
		commentPrivate.GET("/resource/:id", commentController.GetCommentsByResource)	
		commentPrivate.PUT("/update/:id", commentController.DeleteComment)	
	}

	bookmarkPrivate := router.Group("/bookmark")
	bookmarkPrivate.Use(middleware.AuthenticationMiddleware())
	{
		bookmarkPrivate.POST("/create", bookmarkController.CreateBookmark)
		bookmarkPrivate.GET("/:id", bookmarkController.GetBookmarkByID)	
		bookmarkPrivate.DELETE("/delete/:id", bookmarkController.DeleteBookmarkByID)
		bookmarkPrivate.GET("/all", bookmarkController.GetAllBookmarks)
	}

	admin := router.Group("/admin", middleware.AuthenticationMiddleware())
	{
		admin.PUT("/:id", userController.UpdateUser)
		admin.GET("/resources/list", adminController.GetPendingArticles)
	}
}