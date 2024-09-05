package auth

import (
	"goauth/repository"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

type Handler struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
}

func NewHandler(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) *Handler {
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://localhost:5173/callback"),
	)

	return &Handler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
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

	// Generate session token
	sessionToken := generateToken()

	// Store session token in database
	if err := h.tokenRepo.StoreSessionToken(c, user.ID, sessionToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session token"})
		return
	}

	// Generate refresh token
	refreshToken := generateToken()

	// Store refresh token in database
	if err := h.tokenRepo.StoreRefreshToken(c, user.ID, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	// Set session cookie
	c.SetCookie("session", sessionToken, 3600, "/", "", false, false) // 1 hour expiration, not HTTP-only

	// Set refresh cookie
	c.SetCookie("refresh", refreshToken, 30*24*3600, "/", "", false, true) // 30 days expiration, HTTP-only

	c.JSON(http.StatusOK, gin.H{"message": "Authentication successful"})
}

func (h *Handler) HandleRefresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token provided"})
		return
	}

	userID, err := h.tokenRepo.VerifyRefreshToken(c, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new session token
	newSessionToken := generateToken()

	// Store new session token in database
	if err := h.tokenRepo.StoreSessionToken(c, userID, newSessionToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session token"})
		return
	}

	// Generate new refresh token
	newRefreshToken := generateToken()

	// Update refresh token in database
	err = h.tokenRepo.UpdateRefreshToken(c, refreshToken, newRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		return
	}

	// Set new session cookie
	c.SetCookie("session", newSessionToken, 3600, "/", "", false, false) // 1 hour expiration, not HTTP-only

	// Set new refresh cookie
	c.SetCookie("refresh", newRefreshToken, 30*24*3600, "/", "", false, true) // 30 days expiration, HTTP-only

	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed successfully"})
}

func (h *Handler) HandleLogout(c *gin.Context) {
	sessionToken, _ := c.Cookie("session")
	if sessionToken != "" {
		_ = h.tokenRepo.InvalidateSessionToken(c, sessionToken)
	}

	refreshToken, _ := c.Cookie("refresh")
	if refreshToken != "" {
		_ = h.tokenRepo.InvalidateRefreshToken(c, refreshToken)
	}

	// Clear session cookie
	c.SetCookie("session", "", -1, "/", "", false, false)

	// Clear refresh cookie
	c.SetCookie("refresh", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session")
		if err != nil || sessionToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token provided"})
			c.Abort()
			return
		}

		// Verify the session token
		_, err = h.tokenRepo.GetUserIDFromSessionToken(c, sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func generateToken() string {
	return uuid.New().String()
}
