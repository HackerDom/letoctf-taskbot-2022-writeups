package api

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"

	"github.com/HackerDom/letoctf-taskbot-2022-tasks/filestore/pkg/storage"
)

const userSessionKey = "user"
const userIdSessionValueKey = "userId"

type ServerOptions struct {
	Filestore    storage.FileStorage
	DbStore      storage.DbStorage
	SessionStore sessions.Store
	HmacKey      []byte
}

type server struct {
	*http.ServeMux

	filestore    storage.FileStorage
	dbStore      storage.DbStorage
	sessionStore sessions.Store
	hmacKey      []byte
}

func NewServer(
	logger *zerolog.Logger,
	wasmFs http.FileSystem,
	webuiFs http.FileSystem,
	opts *ServerOptions,
) http.Handler {
	s := &server{
		filestore:    opts.Filestore,
		dbStore:      opts.DbStore,
		ServeMux:     &http.ServeMux{},
		sessionStore: opts.SessionStore,
		hmacKey:      opts.HmacKey,
	}

	fsHandlers := map[string]http.Handler{
		"/":      http.FileServer(webuiFs),
		"/wasm/": http.StripPrefix("/wasm", http.FileServer(wasmFs)),
	}
	for pattern, handler := range fsHandlers {
		s.Handle(pattern, handler)
	}

	routes := map[string]apiHandler{
		"/register": s.auth(true),
		"/login":    s.auth(false),
		"/logout":   s.logout,
		"/upload":   s.uploadFile,
		"/list":     s.listFiles,
		"/owner":    s.getFileOwner,
		"/get":      s.getFile,
		"/userid":   s.getCurrentUserId,
	}

	// применяются в прямом порядке
	middlewares := collectMiddlewares(
		loggerMiddleware(logger),
		corsMiddleware,
		authMiddleware(opts.SessionStore),
	)

	for path, handler := range routes {
		h := middlewares(responseMiddleware(handler))
		s.HandleFunc(path, toHttpHandler(context.Background(), h))
	}

	return s
}
