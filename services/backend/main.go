package main

import (
	"fmt"
	"goauth/model"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"goauth/auth"
	"goauth/repository"
)

func main() {
	_ = godotenv.Load(".env")

	// Initialize database
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(&model.User{}, &model.RefreshToken{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// Initialize auth handler
	authHandler, err := auth.NewHandler(userRepo, tokenRepo)
	if err != nil {
		log.Fatalf("Failed to initialize auth handler: %v", err)
	}

	// Set up Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:4173"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Auth routes
	r.GET("/auth/google", authHandler.HandleLogin)
	r.GET("/auth/google/callback", authHandler.HandleCallback)
	r.POST("/auth/refresh", authHandler.HandleRefresh)
	r.GET("/auth/logout", authHandler.HandleLogout)

	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.GET("/user/profile", handleUserProfile())
	}

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}

func handleUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
