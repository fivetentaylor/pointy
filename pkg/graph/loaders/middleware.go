package loaders

import (
	"context"
	"net/http"
	"time"

	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/vikstrous/dataloadgen"
)

// Middleware injects data loaders into the context (for HTTP requests)
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders()
		r = r.WithContext(context.WithValue(r.Context(), LoadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	l, _ := ctx.Value(LoadersKey).(*Loaders)
	if l == nil {
		return NewLoaders()
	}
	return l
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	UserLoader           *dataloadgen.Loader[string, *models.User]
	DocumentOwnerLoader  *dataloadgen.Loader[string, *models.User]
	DocumentAccessLoader *dataloadgen.Loader[DocumentAccessInput, string]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders() *Loaders {
	// define the data loader
	ur := &userReader{}
	dor := &documentOwnerReader{}
	return &Loaders{
		UserLoader:           dataloadgen.NewLoader(ur.getUsers, dataloadgen.WithWait(time.Millisecond)),
		DocumentOwnerLoader:  dataloadgen.NewLoader(dor.getDocumentOwners, dataloadgen.WithWait(time.Millisecond)),
		DocumentAccessLoader: dataloadgen.NewLoader(getDocumentAccesss, dataloadgen.WithWait(time.Millisecond)),
	}
}
