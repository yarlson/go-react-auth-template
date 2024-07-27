package google

import (
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

type AuthProvider struct {
	oauthConfig *oauth2.Config
}

func NewGoogleAuthProvider(clientID string, clientSecret string) *AuthProvider {
	return &AuthProvider{
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *AuthProvider) GetAuthURL() string {
	return g.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (g *AuthProvider) ExchangeCode(code string) (*oauth2.Token, error) {
	return g.oauthConfig.Exchange(context.Background(), code)
}

func (g *AuthProvider) GetUserInfo(token *oauth2.Token) (*User, error) {
	client := g.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo User
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
