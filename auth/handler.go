package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goauth/model"
	"goauth/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"goauth/provider"

	"golang.org/x/oauth2"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type UserIdContextKey struct{}

type UserRepository interface {
	GetOrCreateUser(ctx context.Context, email, firstName, lastName string) (model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID string, refreshToken string) error
	VerifyRefreshToken(ctx context.Context, refreshToken string) (string, error)
	UpdateRefreshToken(ctx context.Context, oldRefreshToken, newRefreshToken string) error
}

type Provider interface {
	GetAuthURL() string
	ExchangeCode(state, code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*provider.User, error)
}

type Handler struct {
	userRepo     UserRepository
	tokenRepo    TokenRepository
	authProvider Provider
	jwtSecret    []byte
}

func NewHandler(userRepo UserRepository, tokenRepo TokenRepository, authProvider Provider, jwtSecret string) *Handler {
	return &Handler{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		authProvider: authProvider,
		jwtSecret:    []byte(jwtSecret),
	}
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.authProvider.GetAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	token, err := h.authProvider.ExchangeCode(state, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	googleUser, err := h.authProvider.GetUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var user model.User
	err = utils.RetryWithBackoff(func() error {
		var err error
		user, err = h.userRepo.GetOrCreateUser(r.Context(), googleUser.Email, googleUser.FirstName, googleUser.LastName)
		return err
	}, 3)
	if err != nil {
		http.Error(w, "Failed to process user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	jwtString, err := generateJWT(user.ID, h.jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshToken := generateRefreshToken()

	// Store refresh token in database
	if err := h.tokenRepo.StoreRefreshToken(r.Context(), user.ID, refreshToken); err != nil {
		http.Error(w, "Failed to store refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token":        jwtString,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := h.tokenRepo.VerifyRefreshToken(r.Context(), request.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate new JWT
	newJWT, err := generateJWT(userID, h.jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	// Generate new refresh token
	newRefreshToken := generateRefreshToken()

	// Update refresh token in database
	err = utils.RetryWithBackoff(func() error {
		return h.tokenRepo.UpdateRefreshToken(r.Context(), request.RefreshToken, newRefreshToken)
	}, 3)
	if err != nil {
		http.Error(w, "Failed to update refresh token", http.StatusInternalServerError)
		return
	}

	// 4. Send both new tokens back to the client
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token":        newJWT,
		"refreshToken": newRefreshToken,
	})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		userID, err := verifyJWT(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdContextKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateJWT(userID string, jwtSecret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiration
	})

	return token.SignedString(jwtSecret)
}

func generateRefreshToken() string {
	return uuid.New().String()
}

func verifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("invalid user ID in token")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}
