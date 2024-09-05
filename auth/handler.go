package auth

import (
	"encoding/json"
	"fmt"
	"goauth/repository"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

type Handler struct {
	userRepo     *repository.UserRepository
	tokenRepo    *repository.TokenRepository
	secureCookie *securecookie.SecureCookie
}

type SessionData struct {
	UserID     string `json:"userId"`
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	PictureURL string `json:"pictureUrl"`
}

func NewHandler(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) (*Handler, error) {
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			"http://localhost:5173/callback",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		),
	)

	hashKey := []byte(os.Getenv("HASH_KEY"))
	blockKey := []byte(os.Getenv("BLOCK_KEY"))

	if len(hashKey) == 0 || len(blockKey) == 0 {
		return nil, fmt.Errorf("HASH_KEY and BLOCK_KEY must be set")
	}

	if len(hashKey) != 32 || len(blockKey) != 32 {
		return nil, fmt.Errorf("HASH_KEY and BLOCK_KEY must be 32 bytes long")
	}

	secureCookie := securecookie.New(hashKey, blockKey)

	return &Handler{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		secureCookie: secureCookie,
	}, nil
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

	user, err := h.userRepo.GetOrCreateUser(c, gothUser.Email, gothUser.FirstName, gothUser.LastName, gothUser.AvatarURL)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user"})
		return
	}

	// Create session data
	sessionData := SessionData{
		UserID:     user.ID.String(),
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		PictureURL: user.PictureURL,
	}

	// Convert session data to JSON
	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal session data"})
		return
	}

	// Encode session data
	encodedSession, err := h.secureCookie.Encode("session", string(jsonData))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode session data"})
		return
	}

	// Generate refresh token
	refreshToken := uuid.New().String()

	// Store refresh token in database
	if err := h.tokenRepo.StoreRefreshToken(c, user.ID, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	encodedRefreshToken, err := h.secureCookie.Encode("refresh", refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode refresh token"})
		return
	}

	// Set session cookie
	c.SetCookie("session", encodedSession, 60, "/", "", true, true) // 1 hour expiration, secure, HTTP-only

	// Set refresh cookie
	c.SetCookie("refresh", encodedRefreshToken, 30*24*3600, "/", "", true, true) // 30 days expiration, secure, HTTP-only

	c.JSON(http.StatusOK, gin.H{"message": "Authentication successful"})
}

func (h *Handler) HandleRefresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token provided"})
		return
	}

	var decodedRefreshToken string
	if err := h.secureCookie.Decode("refresh", refreshToken, &decodedRefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	userID, err := h.tokenRepo.VerifyRefreshToken(c, decodedRefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	user, err := h.userRepo.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Create new session data
	sessionData := SessionData{
		UserID:     user.ID.String(),
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		PictureURL: user.PictureURL,
	}

	// Convert session data to JSON
	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal session data"})
		return
	}

	// Encode new session data
	encodedSession, err := h.secureCookie.Encode("session", string(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode session data"})
		return
	}

	// Generate new refresh token
	newRefreshToken := uuid.New().String()

	// Update refresh token in database
	err = h.tokenRepo.UpdateRefreshToken(c, refreshToken, newRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		return
	}

	encodedRefreshToken, err := h.secureCookie.Encode("refresh", refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode refresh token"})
		return
	}

	// Set new session cookie
	c.SetCookie("session", encodedSession, 3600, "/", "", true, true) // 1 hour expiration, secure, HTTP-only

	// Set new refresh cookie
	c.SetCookie("refresh", encodedRefreshToken, 30*24*3600, "/", "", true, true) // 30 days expiration, secure, HTTP-only

	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed successfully"})
}

func (h *Handler) HandleLogout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh")
	if refreshToken != "" {
		_ = h.tokenRepo.InvalidateRefreshToken(c, refreshToken)
	}

	// Clear session cookie
	c.SetCookie("session", "", -1, "/", "", true, true)

	// Clear refresh cookie
	c.SetCookie("refresh", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionCookie, err := c.Cookie("session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token provided"})
			c.Abort()
			return
		}

		var encodedData string
		if err = h.secureCookie.Decode("session", sessionCookie, &encodedData); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		var sessionData SessionData
		if err = json.Unmarshal([]byte(encodedData), &sessionData); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session data"})
			c.Abort()
			return
		}

		c.Set("user", sessionData)
		c.Next()
	}
}
