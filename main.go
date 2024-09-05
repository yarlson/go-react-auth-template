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
	err = db.AutoMigrate(&model.User{}, &model.RefreshToken{}, &model.SessionToken{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// Initialize auth handler
	authHandler := auth.NewHandler(userRepo, tokenRepo)

	// Set up Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Update this to match your frontend URL
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
		authorized.GET("/user/profile", handleUserProfile(userRepo, tokenRepo))
	}

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}

func handleUserProfile(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Retrieve the user ID associated with the session token
		userID, err := tokenRepo.GetUserIDFromSessionToken(c, sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			return
		}

		user, err := userRepo.GetUserByID(c, userID)
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
