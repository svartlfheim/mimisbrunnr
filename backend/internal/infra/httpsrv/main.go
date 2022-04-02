package httpsrv

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/internal/app/web"
)

type ApiController interface {
	Routes() http.Handler
	RouteGroup() string
}
type openAPIGenerator interface {
	Generate(bool) ([]byte, error)
}

type ServerConfig interface {
	HTTPAPIEnabled() bool
	HTTPStaticServerEnabled() bool
	HTTPFrontendEnabled() bool
	HTTPOpenAPIEnabled() bool
	GetOpenAPIPath() string
	GetHTTPStaticContentPath() string
	GetHTTPPort() string
	GetListenHost() string
}

type Server struct {
	logger         zerolog.Logger
	apiControllers []ApiController
	openAPIGenerator openAPIGenerator
}

// Add the apiVersion to request context
func apiContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiVersion, err := strconv.Atoi(chi.URLParam(r, "apiVersion"))

		if err != nil {
			panic(err)
		}

		ctx := context.WithValue(r.Context(), web.ApiVersionContextKey, apiVersion)
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

// See: https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
func staticFileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func (s *Server) buildRouter(c ServerConfig) http.Handler {
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

	if c.HTTPOpenAPIEnabled() {
		r.Get("/openapi", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// https://github.com/stoplightio/elements#web-component
			body := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>Elements in HTML</title>
    <!-- Embed elements Elements via Web Component -->
    <script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
  </head>
  <body>
      <elements-api
        apiDescriptionUrl="/static/openapi.json?cachebust=` + strconv.Itoa(int(time.Now().Unix())) + `"
        router="hash"
        layout="sidebar"
      />
  </body>
</html>`
			w.Write([]byte(body))
		})
		r.Get(c.GetOpenAPIPath(), func(w http.ResponseWriter, r *http.Request) {
			out, err := s.openAPIGenerator.Generate(false)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(out)
		})
	}
	if c.HTTPStaticServerEnabled() {
		staticFileServer(r, "/static", http.Dir(c.GetHTTPStaticContentPath()))
	}

	if c.HTTPAPIEnabled() {
		r.Route("/api/v{apiVersion:[0-9]+}", func(r chi.Router) {
			r.Use(ensureJSONRequest)
			r.Use(ensureJSONResponse)
			r.Use(apiContext)

			for _, c := range s.apiControllers {
				r.Mount(
					fmt.Sprintf("/%s", c.RouteGroup()),
					c.Routes(),
				)
			}
		})
	}

	if c.HTTPFrontendEnabled() {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
<html>
<head>
	<meta charset="utf-8">
	<title>Mimisbrunnr</title>
	<script src="/static/js/index.js?cachebust=` + strconv.Itoa(int(time.Now().Unix())) + `"></script>
</head>
<body>
	<h1>From CHI</h1>
</body>
</html>
`))
		})
	}

	return r
}

func (s *Server) Start(c ServerConfig) error {
	// Don't need it outside here yet...
	r := s.buildRouter(c)

	listenOn := fmt.Sprintf("%s:%s", c.GetListenHost(), c.GetHTTPPort())

	s.logger.Info().Str("port", c.GetHTTPPort()).Str("host", c.GetListenHost()).Msg("Starting HTTP server")

	return http.ListenAndServe(listenOn, r)
}

func NewServer(logger zerolog.Logger, apiControllers []ApiController, oAPIGen openAPIGenerator) *Server {
	return &Server{
		logger:         logger,
		apiControllers: apiControllers,
		openAPIGenerator: oAPIGen,
	}
}
