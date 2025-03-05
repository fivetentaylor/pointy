package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"golang.org/x/crypto/bcrypt"
)

type EmailSignupRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	AcceptedTerms bool   `json:"acceptedTerms"`
}

type EmailLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Email struct {
	ClientId string
}

func NewEmail() *Email {
	return &Email{}
}

func (e *Email) EmailLogin(w http.ResponseWriter, r *http.Request) {
	var l EmailLoginRequest

	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		log.Error(fmt.Sprintf("Error decoding request body: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userIdent, err := findUserIdentity(r.Context(), l.Email, "email")
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userIdent.PasswordHash), []byte(l.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	signInUserIdent(r.Context(), *userIdent, w, r)
}

func (e *Email) EmailSignup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newUser EmailSignupRequest

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Error(fmt.Sprintf("Error decoding request body: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newUser.Email == "" || newUser.Password == "" || newUser.Name == "" || !newUser.AcceptedTerms {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	userident := userIdentity{
		Email:        newUser.Email,
		Name:         newUser.Name,
		Provider:     "email",
		PasswordHash: string(hashedPassword),
	}

	signInUserIdent(ctx, userident, w, r)
}
