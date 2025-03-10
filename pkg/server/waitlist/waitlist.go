package waitlist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/render"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/fivetentaylor/pointy/pkg/service/email"
	"github.com/fivetentaylor/pointy/pkg/utils"
	"github.com/fivetentaylor/pointy/pkg/views/auth"
)

type WaitlistManager struct {
}

type WaitlistRequest struct {
	Email string `json:"email"`
}

type WaitlistResponse struct {
	Email string `json:"email"`
}

func WaitlistSuccess(w http.ResponseWriter, r *http.Request) {
	auth.Waitlist(os.Getenv("SEGMENT_KEY")).Render(r.Context(), w)
}

func AddToWaitlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newUser WaitlistRequest

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Error(fmt.Sprintf("Error decoding request body: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newUser.Email == "" {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	user, err := CreateWaitlistUser(ctx, newUser)
	if err != nil {
		render.JSON(w, r, WaitlistResponse{Email: ""})
		return
	}

	render.JSON(w, r, WaitlistResponse{Email: user.Email})
}

func CreateWaitlistUser(ctx context.Context, newUser WaitlistRequest) (*models.WaitlistUser, error) {
	q := env.Query(ctx)
	log := env.SLog(ctx)

	normalizedEmail := utils.NormalizeEmail(newUser.Email)

	if ok := utils.IsValidEmail(normalizedEmail); !ok {
		return nil, fmt.Errorf("invalid email")
	}

	var user *models.WaitlistUser
	err := q.Transaction(func(tx *query.Query) error {
		user = &models.WaitlistUser{
			Email: normalizedEmail,
		}

		err := tx.WaitlistUser.Create(user)
		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}

		return nil
	})

	if err != nil {
		if strings.Index(err.Error(), "duplicate key value violates unique constraint") != -1 {
			return nil, errors.New("Email is already registered")
		}
		return nil, err
	}

	err = email.SendWaitlistEmail(
		ctx,
		normalizedEmail,
	)
	if err != nil {
		log.Error("Error sending waitlist email", "error", err)
	}

	log.Info("New waitlist user", "email", normalizedEmail, "event", "new_waitlist_user")

	return user, nil
}
