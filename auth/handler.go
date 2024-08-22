package auth

import (
	"context"
	"goauth/model"
	"goauth/provider"
	"goauth/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

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

func (h *Handler) HandleLogin(c *gin.Context) {
	url := h.authProvider.GetAuthURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) HandleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	token, err := h.authProvider.ExchangeCode(state, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}

	googleUser, err := h.authProvider.GetUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info: " + err.Error()})
		return
	}

	var user model.User
	err = utils.RetryWithBackoff(func() error {
		var err error
		user, err = h.userRepo.GetOrCreateUser(c, googleUser.Email, googleUser.FirstName, googleUser.LastName)
		return err
	}, 3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user: " + err.Error()})
		return
	}

	// Generate JWT
	jwtString, err := generateJWT(user.ID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token: " + err.Error()})
		return
	}

	// Generate refresh token
	refreshToken := generateRefreshToken()

	// Store refresh token in database
	if err := h.tokenRepo.StoreRefreshToken(c, user.ID, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token: " + err.Error()})
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
	newJWT, err := generateJWT(userID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Generate new refresh token
	newRefreshToken := generateRefreshToken()

	// Update refresh token in database
	err = utils.RetryWithBackoff(func() error {
		return h.tokenRepo.UpdateRefreshToken(c, request.RefreshToken, newRefreshToken)
	}, 3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":        newJWT,
		"refreshToken": newRefreshToken,
	})
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
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", jwt.NewValidationError("invalid user ID in token", jwt.ValidationErrorClaimsInvalid)
		}
		return userID, nil
	}

	return "", jwt.NewValidationError("invalid token", jwt.ValidationErrorSignatureInvalid)
}
