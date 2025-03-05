package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/service/payments"
	views "github.com/teamreviso/code/pkg/views/payments"
)

// /payments/checkout
func (s *Server) Checkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	err := r.ParseForm()
	if err != nil {
		log.Error("error parsing form", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	plan := r.FormValue("plan")

	currentUser, err := env.UserClaim(ctx)
	if err != nil {
		log.Error("error getting current user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info("🧾 Checkout", "plan", plan, "userID", currentUser.Id)

	url, err := payments.Checkout(ctx, currentUser.Id, plan)
	if err != nil {
		log.Error("error creating checkout session", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

// /payments/cancel
func (s *Server) CancelCheckout(w http.ResponseWriter, r *http.Request) {
	s.Drafts(w, r)
}

// /payments/failure
func (s *Server) CheckoutFailure(w http.ResponseWriter, r *http.Request) {
	s.Drafts(w, r)
}

// /payments/success
func (s *Server) CheckoutSuccessful(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	_, err := env.UserClaim(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting user claim: %s", err.Error()))
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)

		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Successfull().Render(context.Background(), w)
}

// /payments/status
func (s *Server) CheckoutStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := r.URL.Query().Get("session_id")

	currentUser, err := env.UserClaim(r.Context())
	if err != nil {
		log.Error(fmt.Sprintf("Error getting user claim: %s", err.Error()))
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)

		return
	}

	status, err := payments.Status(ctx, sessionID, currentUser.Id)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting payment status: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bts, err := json.Marshal(status)
	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling payment status: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bts)
}
