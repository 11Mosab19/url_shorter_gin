package main

import (
	"log"
	"url_shorter_gin/database"
	"url_shorter_gin/handlers"
	"url_shorter_gin/middlewares"
	"url_shorter_gin/repository"
	"url_shorter_gin/service"

	"github.com/gin-gonic/gin"
)

func main() {

	//connect to database
	database := &database.Database{}
	err := database.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer database.CloseConnection()

	//repos
	usersRepo := &repository.UserRepository{DB: database}
	urlsRepo := &repository.UrlRepository{DB: database}
	analyticsRepo := &repository.AnalyticsRepository{DB: database}

	//services
	usersService := &service.UserService{Repo: *usersRepo}
	urlsService := &service.UrlService{Repo: *urlsRepo}
	authService := &service.AuthService{US: *usersService}
	analyticsService := &service.AnalyticsService{Repo: *analyticsRepo}

	//handlers
	userHandler := &handlers.UsersHandler{US: *usersService, UrS: *urlsService}
	urlsHandler := &handlers.UrlsHandler{UrS: *urlsService, AH: *analyticsService}
	authHandler := &handlers.AuthHandler{AS: *authService}

	//middlewares
	authMiddleware := &middlewares.AuthenticateMiddleware{AS: *authService}

	r := gin.Default()
	r.Use(middlewares.ErrorMiddleware())
	{
		myApi := r.Group("/api")
		{
			auth := myApi.Group("/auth")
			auth.POST("/login", authHandler.LoginHandler)
			auth.POST("/register", authHandler.RegisterHandler)
		}
		{
			user := myApi.Group("/user")
			user.Use(authMiddleware.AuthMiddleware())
			user.GET("/", userHandler.GetUserProfile)
			user.PATCH("/", userHandler.UpdateUser)
		}
		{
			url := myApi.Group("/url")
			url.Use(authMiddleware.AuthMiddleware())
			url.POST("/", urlsHandler.CreateUrl)
			url.DELETE("/:id", urlsHandler.Delete)
			url.PATCH("/:id", urlsHandler.UpdateUrl)
		}

	}
	r.GET("/:shortCode", urlsHandler.RedirectByShortCode)
	r.POST("/:shortCode/verify", urlsHandler.RedirectVerifiedUrlByShortcode)

	r.Run()
}
