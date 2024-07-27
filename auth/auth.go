package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"goauth/repository"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserIdContextKey struct{}

type UserRepository interface {
	GetOrCreateUser(email string) (repository.User, error)
	GetUserByID(id uint) (repository.User, error)
}

type TokenRepository interface {
	StoreRefreshToken(userID uint, refreshToken string) error
	VerifyRefreshToken(refreshToken string) (uint, error)
	UpdateRefreshToken(userID uint, newRefreshToken string) error
}

type Auth struct {
	googleOauthConfig *oauth2.Config
	userRepo          UserRepository
	tokenRepo         TokenRepository
}

func NewAuth(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *Auth {
	return &Auth{
		googleOauthConfig: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (a *Auth) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := a.googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *Auth) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := a.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := a.googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := a.userRepo.GetOrCreateUser(userInfo.Email)
	if err != nil {
		http.Error(w, "Failed to process user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	jwtString, err := a.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshToken := a.GenerateRefreshToken()

	// Store refresh token in database
	if err := a.tokenRepo.StoreRefreshToken(user.ID, refreshToken); err != nil {
		http.Error(w, "Failed to store refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token":        jwtString,
		"refreshToken": refreshToken,
	})
}

func (a *Auth) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := a.tokenRepo.VerifyRefreshToken(request.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate new JWT
	tokenString, err := a.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Generate new refresh token
	newRefreshToken := a.GenerateRefreshToken()

	// Update refresh token in database
	if err := a.tokenRepo.UpdateRefreshToken(userID, newRefreshToken); err != nil {
		http.Error(w, "Failed to update refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token":        tokenString,
		"refreshToken": newRefreshToken,
	})
}

func (a *Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		userID, err := a.VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdContextKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Auth) GenerateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (a *Auth) GenerateRefreshToken() string {
	return uuid.New().String()
}

func (a *Auth) VerifyJWT(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(float64)
		if !ok {
			return 0, errors.New("invalid user ID in token")
		}
		return uint(userID), nil
	}

	return 0, errors.New("invalid token")
}
