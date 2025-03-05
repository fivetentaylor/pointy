package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teamreviso/code/pkg/admin/templates"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

func GetThreads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")
	usersTbl := env.Query(ctx).User

	log.Info("[admin] loading threads", "docID", docID)

	dydb := env.Dynamo(ctx)
	threads, err := dydb.GetAllThreadsForDoc(docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uniqueUserIds := make(map[string]bool)
	for _, t := range threads {
		uniqueUserIds[t.UserID] = true
	}
	userIds := make([]string, 0, len(uniqueUserIds))
	for k := range uniqueUserIds {
		userIds = append(userIds, k)
	}

	users, err := usersTbl.Where(usersTbl.ID.In(userIds...)).Find()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	threadsByUserId := make(map[string][]*dynamo.Thread)
	for _, t := range threads {
		threadsByUserId[t.UserID] = append(threadsByUserId[t.UserID], t)
	}

	outputMap := make(map[*models.User][]*dynamo.Thread, len(users))
	for _, u := range users {
		outputMap[u] = threadsByUserId[u.ID]
	}

	templates.Threads(docID, outputMap).Render(ctx, w)
}

func GetThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")
	threadID := chi.URLParam(r, "threadId")

	log.Info("[admin] loading thread", "docID", docID, "threadID", threadID)

	dydb := env.Dynamo(ctx)
	thread, err := dydb.GetThreadWithoutUserId(docID, threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messages, err := dydb.GetMessagesForThread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.Thread(docID, thread, messages).Render(ctx, w)
}
