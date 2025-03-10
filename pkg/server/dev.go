package server

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/fivetentaylor/pointy/pkg/service/email/templates"
)

const preheaderExample = "this is the preview shown before the email is opened"

var previewEmails = map[string]func(context.Context, io.Writer) error{
	"magic": func(ctx context.Context, w io.Writer) error {
		return templates.MagicLinkHTML(preheaderExample, "colleen@revi.so", "https://bit.ly/3dV6cFr").Render(ctx, w)
	},
	"waitlist": func(ctx context.Context, w io.Writer) error {
		return templates.WaitlistHTML(preheaderExample, "justin@revi.so").Render(ctx, w)
	},
	"firstopen": func(ctx context.Context, w io.Writer) error {
		doc, err := query.GetRandomDocument(env.Query(ctx))
		if err != nil {
			return err
		}

		user, err := query.GetRandomUser(env.Query(ctx))
		if err != nil {
			return err
		}

		return templates.FirstOpenHTML(preheaderExample, user, doc).Render(ctx, w)
	},
	"mention_share": func(ctx context.Context, w io.Writer) error {
		doc, err := query.GetRandomDocument(env.Query(ctx))
		if err != nil {
			return err
		}

		//user, err := query.GetRandomUser(env.Query(ctx))
		// get user by email
		userTbl := env.Query(ctx).User
		log := env.Log(ctx)
		log.Info("userTbl", "userTbl", userTbl)
		user, err := userTbl.Where(userTbl.Email.Eq("justin@revi.so")).First()
		if err != nil {
			return err
		}

		message := &models.TimelineMessageV1{
			Content:           "hi @@OnVzZXI6ZmQ5ZDE4OWYtYjlkZS00OWE5LThhNTUtNjQyODNhM2U1YzA5OmptcmVpZHkrdDEwMDFAZ21haWwuY29t@@ ",
			ContentAddress:    "[[\"root\",0],[\"q\",1],[[\"root\",0],[\"q\",2],[\"1\",32]]]",
			SelectionStartId:  "",
			SelectionEndId:    "",
			SelectionMarkdown: "yo what's up",
		}

		return templates.MentionShareHTML(preheaderExample, &templates.MentionShareData{
			Document:  doc,
			FromUser:  user,
			Message:   message,
			Selection: "<p>Document text that was commented</p>",
		}).Render(ctx, w)
	},
	"share": func(ctx context.Context, w io.Writer) error {
		return templates.ShareDocHTML("", "Meeting notes", "Taylor", "This doc is awesome", "Join Document", "https://bit.ly/3dV6cFr").Render(ctx, w)
	},
}

func (s *Server) PreviewEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	if preview, ok := previewEmails[email]; ok {
		ctx := r.Context()
		err := preview(s.SES.AttachHostValues(ctx), w)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		return
	}

	w.Write([]byte("Email not found"))
}

func (s *Server) EmailMe(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	if preview, ok := previewEmails[email]; ok {
		ctx := r.Context()
		bts := bytes.NewBuffer([]byte{})
		err := preview(s.SES.AttachHostValues(ctx), bts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		currentUser, err := env.UserClaim(r.Context())
		if err != nil {
			http.Error(w, "Error getting user claim", http.StatusUnauthorized)
			return
		}

		err = s.SES.SendRawEmail("", currentUser.Email, "Test Email", "", bts.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Email sent"))
		return
	}

	w.Write([]byte("Email not found"))
}
