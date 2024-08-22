package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"goauth/auth"
	"goauth/provider/google"
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
	err = db.AutoMigrate(&repository.User{}, &repository.RefreshToken{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// Initialize Google auth provider
	googleAuthProvider := google.NewAuthProvider(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_HOSTED_DOMAIN"),
		os.Getenv("GOOGLE_REDIRECT_URL"),
	)
	defer googleAuthProvider.Stop()

	// Initialize auth
	authHandler := auth.NewHandler(userRepo, tokenRepo, googleAuthProvider, os.Getenv("JWT_SECRET"))

	// Set up Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))

	// Auth routes
	r.GET("/auth/google", gin.WrapF(authHandler.HandleLogin))
	r.GET("/auth/google/callback", gin.WrapF(authHandler.HandleCallback))
	r.POST("/auth/refresh", gin.WrapF(authHandler.HandleRefreshToken))

	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(func(c *gin.Context) {
		authHandler.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	})
	{
		authorized.GET("/user/profile", handleUserProfile(userRepo))
	}

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}

func handleUserProfile(userRepo auth.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Value("userId").(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}

		user, err := userRepo.GetUserByID(c.Request.Context(), userID)
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
