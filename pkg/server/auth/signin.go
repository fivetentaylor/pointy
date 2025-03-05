package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/charmbracelet/log"
	"github.com/go-chi/render"
	"github.com/posthog/posthog-go"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
)

var UserNotAllowedError = fmt.Errorf("user not allowed")

type LoginInResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    string `json:"id"`
}

var emailsList = os.Getenv("EMAIL_ALLOW_LIST")

func signInUser(ctx context.Context, userident userIdentity, w http.ResponseWriter, _ *http.Request) (string, error) {
	/*
	if !isUserAllowed(ctx, userident) && !isDomainAllowed(ctx, userident) {
		waitlist.CreateWaitlistUser(ctx, waitlist.WaitlistRequest{Email: userident.Email})
		return "", UserNotAllowedError
	}
		*/

	token, err := userident.token(ctx)
	if err != nil {
		return "", err
	}

	atCookie := http.Cookie{
		Name:     constants.CookieName,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(UserTokenExpiresIn),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &atCookie)

	return token, nil
}

func signInUserIdent(ctx context.Context, userident userIdentity, w http.ResponseWriter, r *http.Request) {
	token, err := signInUser(ctx, userident, w, r)
	if err != nil {
		log.Errorf("error signing in user: %e", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := userident.User(ctx)
	if err != nil {
		log.Errorf("error loading user: %e", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	render.JSON(w, r, LoginInResponse{
		Token: token,
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser, err := env.UserClaim(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("error loading user claim: %s", err.Error()))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userTbl := env.Query(ctx).User
	user, err := userTbl.Where(userTbl.ID.Eq(currentUser.Id)).First()
	if err != nil {
		log.Error(fmt.Sprintf("error loading user: %s", err.Error()))
	}

	userident := userIdentity{
		Email:    user.Email,
		Name:     user.Name,
		Picture:  user.Picture,
		Provider: user.Provider,
	}

	token, err := userident.token(ctx)
	if err != nil {
		log.Errorf("error refreshing token: %e", err)
		return
	}

	atCookie := http.Cookie{
		Name:     constants.CookieName,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(UserTokenExpiresIn),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &atCookie)

	render.JSON(w, r, map[string]string{"status": "ok"})
}

func isUserAllowed(ctx context.Context, userident userIdentity) bool {
	emails := strings.Split(emailsList, ",")
	waitlistUsersTbl := env.Query(ctx).WaitlistUser
	existingUser, err := waitlistUsersTbl.Where(waitlistUsersTbl.Email.Eq(userident.Email)).First()
	if err != nil {
		log.Infof("Error checking if user is on waitlist: %s", err.Error())
	}
	return slices.Contains(emails, userident.Email) || (existingUser != nil && existingUser.AllowAccess == true)
}

func isDomainAllowed(ctx context.Context, userident userIdentity) bool {
	log.Infof("checking if domain is allowed for user: %s", userident.Email)
	phClient := env.Posthog(ctx)
	if phClient == nil {
		return false
	}
	isMyFlagEnabled, err := phClient.IsFeatureEnabled(
		posthog.FeatureFlagPayload{
			Key:              "is-domain-allowed",
			DistinctId:       userident.Email,
			PersonProperties: posthog.NewProperties().Set("email", userident.Email),
		})

	if err != nil {
		log.Errorf("error checking if feature flag is enabled: %s", err)
		return false
	}

	log.Infof("isMyFlagEnabled: %t for user: %s", isMyFlagEnabled, userident.Email)

	return isMyFlagEnabled == true
}
