package admin

import (
	"net/http"

	"github.com/teamreviso/code/pkg/admin/templates"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/service/prompts"
)

func GetPrompts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tbl := env.Query(ctx).Prompt

	prompts, err := tbl.Order(tbl.ID).Find()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.Prompts(prompts).Render(ctx, w)
}

func RefreshPrompts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := prompts.Refresh(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/prompts", http.StatusSeeOther)
}
