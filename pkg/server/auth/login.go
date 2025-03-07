package auth

import (
	"context"
	"net/http"
	"os"

	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/server/auth/types"
	"github.com/fivetentaylor/pointy/pkg/views/auth"
	"github.com/fivetentaylor/pointy/pkg/views/utils"
)

func (m *Manager) GetLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.SLog(ctx)

	queryParams := r.URL.Query()

	next := queryParams.Get("next")
	sb := queryParams.Get("sb")

	state := types.State{
		Next:    next,
		Sidebar: sb,
	}

	currentUser, err := env.UserClaim(ctx)
	if err == nil && currentUser != nil {
		log.Info("login", "userID", currentUser.Id, "event", "user_login")
		nextURL := stateToNextURL(&state)
		if nextURL != "" {
			http.Redirect(w, r, nextURL, http.StatusFound)
		} else {
			http.Redirect(w, r, "/drafts", http.StatusFound)
		}
		return
	}

	stateString, err := m.generateStateString(state)
	if err != nil {
		log.Error("error generating state", "error", err)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		utils.Error().Render(context.Background(), w)
		return
	}

	auth.Login(os.Getenv("SEGMENT_KEY"), stateString, state).Render(r.Context(), w)
}
