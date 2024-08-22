package main

import (
	"fmt"
	"goauth/model"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"goauth/auth"
	"goauth/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	authHandler := auth.NewHandler(userRepo, tokenRepo, os.Getenv("JWT_SECRET"))

	// Set up Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))

	// Auth routes
	r.GET("/auth/google", authHandler.HandleLogin)
	r.GET("/auth/google/callback", authHandler.HandleCallback)
	r.POST("/auth/refresh", authHandler.HandleRefreshToken)
	r.GET("/auth/logout", authHandler.HandleLogout)

	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.GET("/user/profile", handleUserProfile(userRepo))
	}

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}

func handleUserProfile(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user, err := userRepo.GetUserByID(c, uuid.MustParse(userID.(string)))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		})
	}
}
