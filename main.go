package main

import (
	"api-service/config"
	"api-service/controllers"
	"api-service/db"
	"api-service/repositories"
	"api-service/services"
	"api-service/migrations"
	"api-service/routers"
	"log"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Load configuration
    cfg := config.LoadConfig()

    // Connect to the database
    database, err := db.ConnectDatabase(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Migrate model
    if err := migrations.MigrateModels(database); err != nil {
        log.Fatalf("Database migration failed: %v", err)
    }

    // Create router
    router := gin.Default()

	// Custom CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
    // Initialize dependencies

	authRepo := repositories.NewAuthRepository(database)
	authService := services.NewAuthService(*authRepo)
	
    userRepo := repositories.NewUserRepository(database)
    userService := services.NewUserService(*userRepo)
    
    commentRepo := repositories.NewCommentRepository(database)
    commentService := services.NewCommentService(*commentRepo)
    
    bookmarkRepo := repositories.NewBookmarkRepository(database)
	bookmarkService := services.NewBookmarkService(*bookmarkRepo)

    reportRepo := repositories.NewReportRepository(database)
    reportService := services.NewReportService(*reportRepo)

    resourceRep := repositories.NewResourceRepository(database)
    resourceService := services.NewResourceService(*resourceRep)

    adminRepo := repositories.NewAdminRepository(database)
    adminService := services.NewAdminService(*adminRepo)

    resourceController := controllers.NewResourceController(*&resourceService)
    reportController := controllers.NewReportController(*reportService)
    authController := controllers.NewAuthController(*authService)
    userController := controllers.NewUserController(*userService)
	bookmarkController := controllers.NewBookmarkController(*bookmarkService)
    commentController := controllers.NewCommentController(*commentService)
    adminController := controllers.NewAdminController(*adminService)
    
    // Register all routes by passing the router and dependencies
    routers.RegisterRoutes(router, authController, userController, 
        bookmarkController, commentController, 
        reportController, resourceController,
        adminController) 

    // Start the server
    if err := router.Run(":7099"); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
    gin.SetMode(gin.ReleaseMode)
}