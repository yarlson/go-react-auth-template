package google

import (
	"context"
	"encoding/json"
	"goauth/provider"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthProvider struct {
	oauthConfig  *oauth2.Config
	hostedDomain string
}

func NewAuthProvider(clientID, clientSecret, hostedDomain string) *AuthProvider {
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
		hostedDomain: hostedDomain,
	}
}

func (g *AuthProvider) GetAuthURL() string {
	opts := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}
	if g.hostedDomain != "" {
		opts = append(opts, oauth2.SetAuthURLParam("hd", g.hostedDomain))
	}
	return g.oauthConfig.AuthCodeURL("state", opts...)
}

func (g *AuthProvider) ExchangeCode(code string) (*oauth2.Token, error) {
	return g.oauthConfig.Exchange(context.Background(), code)
}

func (g *AuthProvider) GetUserInfo(token *oauth2.Token) (*provider.User, error) {
	client := g.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo provider.User
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
