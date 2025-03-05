package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/server/auth/types"
	"github.com/teamreviso/code/pkg/service/email"
	"github.com/teamreviso/code/pkg/views/auth"
	authView "github.com/teamreviso/code/pkg/views/auth"
	viewUtils "github.com/teamreviso/code/pkg/views/utils"
)

var OneTimeAccessLinkExpireIn = time.Minute * 15

type AccessLinkRequest struct {
	Token string `json:"token"`
}

type SendAccessLinkRequest struct {
	Email string  `json:"email"`
	Next  *string `json:"next"`
}

func SendAccessLink(ctx context.Context, emailaddr, next string) error {
	userTbl := env.Query(ctx).User

	user, err := userTbl.
		Where(userTbl.Email.Eq(emailaddr)).
		First()

	if err != nil {
		if err.Error() == "record not found" {
			userIdentity := &userIdentity{
				Email:    emailaddr,
				Name:     emailaddr,
				Provider: "accesslink",
			}

			user, err = userIdentity.FindOrCreateUser(ctx)
		} else {
			return err
		}
	}

	oneTimeAccessLinkTbl := env.Query(ctx).OneTimeAccessToken

	// Use all existing links
	_, err = oneTimeAccessLinkTbl.
		Where(oneTimeAccessLinkTbl.UserID.Eq(user.ID)).
		Update(oneTimeAccessLinkTbl.IsUsed, true)
	if err != nil {
		return err
	}

	newOneTimeAccessLink := models.OneTimeAccessToken{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(OneTimeAccessLinkExpireIn),
		IsUsed:    false,
	}
	err = oneTimeAccessLinkTbl.Create(&newOneTimeAccessLink)
	if err != nil {
		return err
	}

	err = email.SendMagicLinkEmail(
		ctx,
		emailaddr,
		newOneTimeAccessLink.Token,
		next,
	)

	return nil
}

func GetMagicLink(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	next := r.URL.Query().Get("next")
	log := env.SLog(r.Context())

	userIdent, err := findUserIdentityByOneTimeAccessToken(r.Context(), code)
	if err != nil {
		log.Error(fmt.Sprintf("Error finding user identity: %s", err.Error()))
		auth.MagicFailed(os.Getenv("SEGMENT_KEY"), "Your access link may have expired.").Render(r.Context(), w)
		return
	}

	_, err = signInUser(r.Context(), *userIdent, w, r)
	if err != nil {
		if err == UserNotAllowedError {
			http.Redirect(w, r, "/waitlist/success", http.StatusTemporaryRedirect)
			return
		}

		log.Error(fmt.Sprintf("Error signing in user: %s", err.Error()))
		auth.MagicFailed(os.Getenv("SEGMENT_KEY"), "Please reach out to contact@revi.so if this continues.").Render(r.Context(), w)
		return
	}

	if next != "" && next[0] == '/' {
		http.Redirect(w, r, next, http.StatusFound)
		return
	}

	log.Info("login", "userID", userIdent.ID, "event", "user_login", "event_type", "magic_link_login")
	http.Redirect(w, r, "/drafts", http.StatusFound)
}

func SendMagicLink(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		authView.LoginForm("", types.State{}, "Invalid form").Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	nextInput := r.FormValue("next")
	sbInput := r.FormValue("sb")

	var next string
	if nextInput != "" && nextInput[0] == '/' {
		next = nextInput
		if sbInput != "" {
			next += "?sb=" + sbInput
		}
	}

	err := SendAccessLink(r.Context(), email, next)
	if err != nil {
		authView.LoginForm(email, types.State{}, "Failed sending access link").Render(r.Context(), w)
		return
	}

	auth.EmailSent(email).Render(r.Context(), w)
}

func sendOneTimeAccessLinkForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		authView.LoginForm("", types.State{}, "Invalid form").Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")

	err := SendAccessLink(r.Context(), email, "")
	if err != nil {
		authView.LoginForm(email, types.State{}, "Failed sending access link").Render(r.Context(), w)
		return
	}

	// TODO - redirect to success page
	viewUtils.Redirect("/").Render(r.Context(), w)
	return
}

func OneTimeAccessLink(w http.ResponseWriter, r *http.Request) {
	var req AccessLinkRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error(fmt.Sprintf("Error decoding request body: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userIdent, err := findUserIdentityByOneTimeAccessToken(r.Context(), req.Token)
	if err != nil {
		log.Error(fmt.Sprintf("Error finding user identity: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	signInUserIdent(r.Context(), *userIdent, w, r)
}
