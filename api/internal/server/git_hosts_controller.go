package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GitHostsController struct {}

func gitHostContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gitHostID := chi.URLParam(r, "gitHostID")

		ctx := context.WithValue(r.Context(), "gitHostID", gitHostID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *GitHostsController) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", c.List)
	r.Get("/search", c.Search)

	r.Post("/", c.Create)

	r.Route("/{gitHostID}", func(r chi.Router) {
		r.Use(gitHostContext)
		r.Get("/", c.Get)
		r.Put("/", c.Update)
		r.Delete("/", c.Delete)
	})

	return r
}

func (c *GitHostsController) RouteGroup() string {
	return "git-hosts"
}

func (c *GitHostsController) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list git hosts"))
}

func (c *GitHostsController) Search(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("search git hosts"))
}

func (c *GitHostsController) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create git host"))
}

func (c *GitHostsController) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get git host: " + r.Context().Value("gitHostID").(string)))
}

func (c *GitHostsController) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update git host: " + r.Context().Value("gitHostID").(string)))
}

func (c *GitHostsController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete git host: " + r.Context().Value("gitHostID").(string)))
}

func NewGitHostsController() *GitHostsController {
	return &GitHostsController{}
}