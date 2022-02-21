package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type Controller interface {
	Routes() http.Handler
	RouteGroup() string
}
type serverConfig interface {
	GetHTTPPort() string
	GetListenHost() string
}

type Server struct {
	logger zerolog.Logger
	controllers []Controller
}

func defaultHandler(n string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("route (%s): %s", r.Context().Value("apiVersion").(string), n)))
	}
}

// Add the apiVersion to request context
func apiContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiVersion := chi.URLParam(r, "apiVersion")

		ctx := context.WithValue(r.Context(), "apiVersion", apiVersion)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) buildRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.GetHead)
	r.Use(middleware.Timeout(60 * time.Second))
  
	r.Route("/api/v{apiVersion:[0-9]+}", func(r chi.Router) {
		r.Use(apiContext)

		for _, c := range(s.controllers) {
			r.Mount(
				fmt.Sprintf("/%s", c.RouteGroup()),
				c.Routes(),
			)
		}
	})

	return r
}

func (s *Server) Start(c serverConfig) error {
	// Don't need it outside here yet...
	r := s.buildRouter()

	listenOn := fmt.Sprintf("%s:%s", c.GetListenHost(), c.GetHTTPPort())

	s.logger.Info().Str("port", c.GetHTTPPort()).Str("host", c.GetListenHost()).Msg("Starting HTTP server")

	return http.ListenAndServe(listenOn, r)
}

func NewServer(logger zerolog.Logger, controllers []Controller) *Server {
	return &Server{
		logger: logger,
		controllers: controllers,
	}
}