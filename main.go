package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"goauth/auth"
	"goauth/provider/google"
	"goauth/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
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
	userRepo := repository.NewGormUserRepository(db)
	tokenRepo := repository.NewGormTokenRepository(db)

	// Initialize Google auth provider
	googleAuthProvider := google.NewAuthProvider(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_HOSTED_DOMAIN"),
		os.Getenv("GOOGLE_REDIRECT_URL"),
	)

	// Initialize auth
	authHandler := auth.NewAuth(userRepo, tokenRepo, googleAuthProvider, os.Getenv("JWT_SECRET"))

	// Set up chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsMiddleware.Handler)

	r.Get("/auth/google", authHandler.HandleLogin)
	r.Get("/auth/google/callback", authHandler.HandleCallback)

	r.Post("/auth/refresh", authHandler.HandleRefreshToken)

	r.Group(func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)
		r.Get("/api/user/profile", handleUserProfile(userRepo))
	})

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleUserProfile(userRepo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.UserIdContextKey{}).(uint)
		user, err := userRepo.GetUserByID(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		})
	}
}
