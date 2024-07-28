package google

import (
	"context"
	"encoding/json"
	"errors"
	"goauth/provider"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthProvider struct {
	oauthConfig   *oauth2.Config
	hostedDomain  string
	stateStore    map[string]time.Time
	stateStoreMux sync.Mutex
	cleanupTicker *time.Ticker
	done          chan bool
}

func NewAuthProvider(clientID, clientSecret, hostedDomain, redirectURL string) *AuthProvider {
	authProvider := &AuthProvider{
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		hostedDomain:  hostedDomain,
		stateStore:    make(map[string]time.Time),
		cleanupTicker: time.NewTicker(5 * time.Minute),
		done:          make(chan bool),
	}

	go authProvider.runCleanup()

	return authProvider
}

func (g *AuthProvider) runCleanup() {
	for {
		select {
		case <-g.done:
			return
		case <-g.cleanupTicker.C:
			g.cleanupExpiredStates()
		}
	}
}

func (g *AuthProvider) Stop() {
	g.cleanupTicker.Stop()
	g.done <- true
}

func (g *AuthProvider) generateState() string {
	state := uuid.New().String()
	g.stateStoreMux.Lock()
	defer g.stateStoreMux.Unlock()
	g.stateStore[state] = time.Now().Add(30 * time.Minute)
	return state
}

func (g *AuthProvider) validateState(state string) bool {
	g.stateStoreMux.Lock()
	defer g.stateStoreMux.Unlock()

	expirationTime, exists := g.stateStore[state]
	if !exists {
		return false
	}
	if time.Now().After(expirationTime) {
		delete(g.stateStore, state)
		return false
	}

	delete(g.stateStore, state)
	return true
}

func (g *AuthProvider) cleanupExpiredStates() {
	g.stateStoreMux.Lock()
	defer g.stateStoreMux.Unlock()
	now := time.Now()
	for state, expirationTime := range g.stateStore {
		if now.After(expirationTime) {
			delete(g.stateStore, state)
		}
	}
}

func (g *AuthProvider) GetAuthURL() string {
	state := g.generateState()
	opts := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}
	if g.hostedDomain != "" {
		opts = append(opts, oauth2.SetAuthURLParam("hd", g.hostedDomain))
	}
	return g.oauthConfig.AuthCodeURL(state, opts...)
}

func (g *AuthProvider) ExchangeCode(state string, code string) (*oauth2.Token, error) {
	if !g.validateState(state) {
		return nil, errors.New("invalid state parameter")
	}

	token, err := g.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	return token, nil
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
