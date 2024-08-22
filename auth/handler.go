package auth

import (
	"context"
	"errors"
	"fmt"
	"goauth/model"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

type UserRepository interface {
	GetOrCreateUser(ctx context.Context, email, firstName, lastName string) (model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error)
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error
	VerifyRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error)
	UpdateRefreshToken(ctx context.Context, oldRefreshToken, newRefreshToken string) error
}

type Handler struct {
	userRepo  UserRepository
	tokenRepo TokenRepository
	jwtSecret []byte
}

func NewHandler(userRepo UserRepository, tokenRepo TokenRepository, jwtSecret string) *Handler {
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://localhost:5173/callback"),
	)

	return &Handler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (h *Handler) HandleLogin(c *gin.Context) {
	provider, err := goth.GetProvider("google")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get provider"})
		return
	}

	state := uuid.New().String()
	sess, err := provider.BeginAuth(state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin authentication"})
		return
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get auth URL"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) HandleCallback(c *gin.Context) {
	provider, err := goth.GetProvider("google")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get provider"})
		return
	}

	state := c.Query("state")

	sess, err := provider.BeginAuth(state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin authentication"})
		return
	}

	_, err = sess.Authorize(provider, c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authorize"})
		return
	}

	gothUser, err := provider.FetchUser(sess)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	user, err := h.userRepo.GetOrCreateUser(c, gothUser.Email, gothUser.FirstName, gothUser.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user"})
		return
	}

	// Generate JWT
	jwtString, err := generateJWT(user.ID.String(), h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Generate refresh token
	refreshToken := generateRefreshToken()

	// Store refresh token in database
	if err := h.tokenRepo.StoreRefreshToken(c, user.ID, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":        jwtString,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) HandleRefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, err := h.tokenRepo.VerifyRefreshToken(c, request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new JWT
	newJWT, err := generateJWT(userID.String(), h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Generate new refresh token
	newRefreshToken := generateRefreshToken()

	// Update refresh token in database
	err = h.tokenRepo.UpdateRefreshToken(c, request.RefreshToken, newRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":        newJWT,
		"refreshToken": newRefreshToken,
	})
}

func (h *Handler) HandleLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		userID, err := verifyJWT(tokenString, h.jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
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

func verifyJWT(tokenString string, jwtSecret []byte) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
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
