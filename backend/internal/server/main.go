package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

type Controller interface {
	Routes() http.Handler
	RouteGroup() string
}
type ServerConfig interface {
	GetHTTPPort() string
	GetListenHost() string
}

type Server struct {
	logger      zerolog.Logger
	controllers []Controller
}

// Add the apiVersion to request context
func apiContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiVersion := chi.URLParam(r, "apiVersion")

		ctx := context.WithValue(r.Context(), apiVersionContextKey, apiVersion)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ensureJSONResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

// See: https://github.com/go-chi/chi/blob/master/middleware/content_type.go
// Basically a cheap rip off
// The official one skipped the check if the body was empty
// This caused 500 errors when we parse the body
func ensureJSONRequest(next http.Handler) http.Handler {
	allowedContentTypes := map[string]bool{
		"application/json": true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
		if i := strings.Index(s, ";"); i > -1 {
			s = s[0:i]
		}

		if _, ok := allowedContentTypes[s]; ok {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnsupportedMediaType)
	})
}

func (s *Server) buildRouter() http.Handler {
	r := chi.NewRouter()
	logger := httplog.NewLogger("server", httplog.Options{
		JSON: true,
	})
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.GetHead)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v{apiVersion:[0-9]+}", func(r chi.Router) {
		r.Use(ensureJSONRequest)
		r.Use(ensureJSONResponse)
		r.Use(apiContext)

		for _, c := range s.controllers {
			r.Mount(
				fmt.Sprintf("/%s", c.RouteGroup()),
				c.Routes(),
			)
		}
	})

	return r
}

func (s *Server) Start(c ServerConfig) error {
	// Don't need it outside here yet...
	r := s.buildRouter()

	listenOn := fmt.Sprintf("%s:%s", c.GetListenHost(), c.GetHTTPPort())

	s.logger.Info().Str("port", c.GetHTTPPort()).Str("host", c.GetListenHost()).Msg("Starting HTTP server")

	return http.ListenAndServe(listenOn, r)
}

func NewServer(logger zerolog.Logger, controllers []Controller) *Server {
	return &Server{
		logger:      logger,
		controllers: controllers,
	}
}
