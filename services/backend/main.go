package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"backend/auth"
	"backend/model"
	"backend/repository"
)

func main() {
	_ = godotenv.Load(".env")

	// Initialize database
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = "data/app.db" // Default path if env var is not set
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{}) // Change this line
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
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:4173", "http://localhost"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes
	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.GET("/google", authHandler.HandleLogin)
		authGroup.GET("/google/callback", authHandler.HandleCallback)
		authGroup.POST("/refresh", authHandler.HandleRefresh)
		authGroup.GET("/logout", authHandler.HandleLogout)
	}

	// Protected routes
	authorized := r.Group("/api/v1")
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
