package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/config"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/views/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOauth struct {
	config *oauth2.Config
}

func NewGoogle(cfg *config.GoogleOauth) *GoogleOauth {
	log.Debug(fmt.Sprintf("Creating google oauth config: %v", cfg))
	config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	return &GoogleOauth{
		config: config,
	}
}

func (m *Manager) GoogleRedirect(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	url := m.google.config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (m *Manager) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.SLog(ctx)

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	authState, valid, err := m.verifyStateString(state)
	if !valid {
		log.Error(fmt.Sprintf("Error state not valid: %s", state))
		auth.AuthFailure(os.Getenv("SEGMENT_KEY"), "Something went wrong").Render(r.Context(), w)
		return
	}

	if err != nil {
		log.Error(fmt.Sprintf("Error verifying state string: %s", err.Error()))
		auth.AuthFailure(os.Getenv("SEGMENT_KEY"), "Something went wrong").Render(r.Context(), w)
		return
	}

	token, err := m.google.config.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Could not get token: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	userInfo, err := FetchGoogleUserInfo(token)
	if err != nil {
		log.Error(fmt.Sprintf("Error fetching user info: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Debug(fmt.Sprintf("User info: %v", userInfo))
	userident := userIdentity{
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		Picture:  &userInfo.Picture,
		Provider: "google",
	}

	/*
	if !isUserAllowed(ctx, userident) && !isDomainAllowed(ctx, userident) {
		waitlist.CreateWaitlistUser(ctx, waitlist.WaitlistRequest{Email: userident.Email})
		http.Redirect(w, r, "/waitlist/success", http.StatusTemporaryRedirect)
		return
	}
		*/

	revisoToken, err := userident.token(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting token: %s", err.Error()))
		auth.AuthFailure(os.Getenv("SEGMENT_KEY"), "Something went wrong").Render(r.Context(), w)
		return
	}

	atCookie := http.Cookie{
		Name:     constants.CookieName,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		Value:    revisoToken,
		Path:     "/",
		Expires:  time.Now().Add(UserTokenExpiresIn),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(w, &atCookie)

	next := stateToNextURL(authState)
	log.Info("next outside redirect", "next", next)

	if next != "" {
		log.Info("inside redirecting to next", "next", next)
		http.Redirect(w, r, next, http.StatusTemporaryRedirect)
		return
	}

	log.Info("login", "userID", userident.ID, "event", "user_login", "event_type", "google_login")
	http.Redirect(w, r, "/drafts", http.StatusTemporaryRedirect)
}

type UserInfo struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

func FetchGoogleUserInfo(token *oauth2.Token) (*UserInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debug(fmt.Sprintf("Google user info: %s", string(data)))

	// Unmarshal the data into UserInfo struct
	var userInfo UserInfo
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}
