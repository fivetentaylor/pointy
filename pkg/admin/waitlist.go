package admin

import (
	"context"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/teamreviso/code/pkg/admin/templates"
	"github.com/teamreviso/code/pkg/client"
	"github.com/teamreviso/code/pkg/env"
)

func WaitlistHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	entries, err := GetWaitlistEntries(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch waitlist entries", http.StatusInternalServerError)
		return
	}

	err = templates.Waitlist(entries).Render(ctx, w)
	if err != nil {
		http.Error(w, "Failed to render waitlist template", http.StatusInternalServerError)
		return
	}
}

// New handler for updating waitlist user access
func UpdateWaitlistAccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	err := UpdateWaitlistAccess(ctx, email)
	if err != nil {
		http.Error(w, "Failed to update waitlist access", http.StatusInternalServerError)
		return
	}

	loopsClient, err := client.NewLoopsClientFromEnv()
	if err != nil {
		log.Info("Failed to create loops client", "error", err)
		http.Error(w, "Failed to create loops client", http.StatusInternalServerError)
		return
	}

	err = loopsClient.SendEvent(ctx, "off-waitlist", client.LoopsContactProperties{Email: email}, nil)
	if err != nil {
		http.Error(w, "Failed to send event to loops", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/waitlist", http.StatusSeeOther)
}

func UpdateWaitlistAccess(ctx context.Context, email string) error {
	waitlistTbl := env.Query(ctx).WaitlistUser
	_, err := waitlistTbl.Where(waitlistTbl.Email.Eq(email)).Update(waitlistTbl.AllowAccess, true)
	return err
}

func GetWaitlistEntries(ctx context.Context) ([]templates.WaitlistEntry, error) {
	waitlistTbl := env.Query(ctx).WaitlistUser
	rows, err := waitlistTbl.Order(waitlistTbl.CreatedAt).Find()
	if err != nil {
		return nil, err
	}

	entries := make([]templates.WaitlistEntry, len(rows))
	for i, row := range rows {
		entries[i] = templates.WaitlistEntry{
			Email:       row.Email,
			CreatedAt:   row.CreatedAt.String(),
			UpdatedAt:   row.UpdatedAt.String(),
			AllowAccess: row.AllowAccess,
		}
	}

	return entries, nil
}
