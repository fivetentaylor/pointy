package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/server/auth/types"
	"github.com/fivetentaylor/pointy/pkg/service/sharing"
	"github.com/fivetentaylor/pointy/pkg/views/auth"
)

func (m *Manager) GetInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := chi.URLParam(r, "code")
	log := env.Log(ctx)

	shareLinkTbl := env.Query(ctx).SharedDocumentLink
	userTbl := env.Query(ctx).User

	currentUser, _ := env.UserClaim(ctx)
	if currentUser != nil {
		log.Infof("User already logged in: %s, joining invite", currentUser.Email)
		m.JoinInvite(w, r)
		return
	}

	sl, err := shareLinkTbl.
		Where(shareLinkTbl.InviteLink.Eq(code)).
		First()
	if err != nil {
		log.Error(fmt.Sprintf("[invite] Error finding share link: %s", err.Error()))
		auth.InviteFailed(os.Getenv("SEGMENT_KEY")).Render(r.Context(), w)
		return
	}

	invitedBy, err := userTbl.
		Where(userTbl.ID.Eq(sl.InviterID)).
		First()
	if err != nil {
		log.Error(fmt.Sprintf("[invite] Error finding user: %s", err.Error()))
		auth.InviteFailed(os.Getenv("SEGMENT_KEY")).Render(r.Context(), w)
		return
	}

	next := fmt.Sprintf("/invite/%s", sl.InviteLink)
	state := types.State{
		Next: next,
	}

	stateString, err := m.generateStateString(state)
	if err != nil {
		log.Error(fmt.Errorf("[invite] error generating state: %w", err))
		auth.InviteFailed(os.Getenv("SEGMENT_KEY")).Render(r.Context(), w)
		return
	}

	auth.Invite(os.Getenv("SEGMENT_KEY"), stateString, state, sl, invitedBy).Render(r.Context(), w)
}

func (m *Manager) JoinInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := chi.URLParam(r, "code")
	currentUser, err := env.UserClaim(ctx)
	if err != nil {
		log.Errorf("[invite] error getting current user: %s", err)
		auth.InviteFailed(os.Getenv("SEGMENT_KEY")).Render(r.Context(), w)
		return
	}

	shareLinkTbl := env.Query(ctx).SharedDocumentLink
	sl, err := shareLinkTbl.
		Where(shareLinkTbl.InviteLink.Eq(code)).
		Where(shareLinkTbl.InviteeEmail.Eq(currentUser.Email)).
		First()
	if err != nil {
		log.Errorf("[invite] error getting share link: %s", err)
		auth.InviteFailed(os.Getenv("SEGMENT_KEY")).Render(r.Context(), w)
		return
	}

	doc, err := sharing.JoinDoc(ctx, sl.DocumentID, currentUser.Id, "write")
	if err != nil {
		log.Errorf("[invite] error joining share link: %s", err)
		auth.InviteFailed(os.Getenv("SEGMENT_KEY")).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/drafts/%s", doc.ID), http.StatusFound)
}
