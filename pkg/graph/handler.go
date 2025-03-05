package graph

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"

	"github.com/teamreviso/code/pkg/graph/loaders"
)

func NewHandler() http.Handler {
	userSrv := handler.New(
		NewExecutableSchema(
			Config{Resolvers: &Resolver{}},
		),
	)
	userSrv.AddTransport(transport.POST{})
	userSrv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin:     func(r *http.Request) bool { return true },
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})
	userSrv.AddTransport(&transport.MultipartForm{
		MaxMemory:     32 << 20,
		MaxUploadSize: 50 << 20,
	})
	userSrv.SetQueryCache(lru.New(1000))
	userSrv.Use(extension.Introspection{})

	wDataLoaders := loaders.Middleware(userSrv)

	return wDataLoaders
}
