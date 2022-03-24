package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/projects/v1"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
)

type ProjectsController interface {
	AddV1(dto v1.AddProjectDTO) commandresult.Result
	ListV1(dto v1.ListProjectsDTO) commandresult.Result
	GetV1(id string) commandresult.Result
	UpdateV1(id string, dto v1.UpdateProjectDTO) commandresult.Result
	DeleteV1(id string) commandresult.Result
}

type ProjectsHandler struct {
	logger           zerolog.Logger
	controller       ProjectsController
	jsonUnmarshaller jsonUnmarshaller
}

func projectContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID := chi.URLParam(r, "projectID")

		ctx := context.WithValue(r.Context(), projectIDContextKey, projectID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func pageContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageID := chi.URLParam(r, "pageID")

		ctx := context.WithValue(r.Context(), pageIDContextKey, pageID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *ProjectsHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", c.List)

	r.Post("/", c.Create)

	r.Route("/{projectID}", func(r chi.Router) {
		r.Use(projectContext)
		r.Get("/", c.Get)
		r.Patch("/", c.Update)
		r.Delete("/", c.Delete)

		r.Route("/pages", func(r chi.Router) {
			r.Get("/", c.ListPages)

			r.Get("/search", c.SearchPages)

			r.Route("/{pageID}", func(r chi.Router) {
				r.Use(pageContext)
				r.Get("/", c.GetPage)
			})
		})
	})
	return r
}

func (c *ProjectsHandler) RouteGroup() string {
	return "projects"
}

func (c *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	dto := v1.ListProjectsDTO{}

	if page != "" {
		if pageAsInt, err := strconv.Atoi(page); err == nil {
			dto.Page = &pageAsInt
		}
	}

	if limit != "" {
		if limitAsInt, err := strconv.Atoi(limit); err == nil {
			dto.Limit = &limitAsInt
		}
	}

	serveResponseForResult(
		c.controller.ListV1(dto),
		c.logger,
		w,
	)
}

func (c *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	dto := v1.AddProjectDTO{}

	if err := c.jsonUnmarshaller.Unmarshal(r, &dto); handleError(w, c.logger, err) {
		return
	}

	serveResponseForResult(
		c.controller.AddV1(dto),
		c.logger,
		w,
	)
}

func (c *ProjectsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(projectIDContextKey).(string)

	serveResponseForResult(
		c.controller.GetV1(id),
		c.logger,
		w,
	)
}

func (c *ProjectsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(projectIDContextKey).(string)

	dto := v1.UpdateProjectDTO{}

	if err := c.jsonUnmarshaller.Unmarshal(r, &dto); handleError(w, c.logger, err) {
		return
	}

	serveResponseForResult(
		c.controller.UpdateV1(id, dto),
		c.logger,
		w,
	)
}

func (c *ProjectsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(projectIDContextKey).(string)

	serveResponseForResult(
		c.controller.DeleteV1(id),
		c.logger,
		w,
	)
}

func (c *ProjectsHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf(
		"get project (%s) page: %s",
		r.Context().Value(projectIDContextKey).(string),
		r.Context().Value(pageIDContextKey).(string),
	)))

	if err != nil {
		panic(err)
	}
}

func (c *ProjectsHandler) SearchPages(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf(
		"search project (%s) pages",
		r.Context().Value(projectIDContextKey).(string),
	)))

	if err != nil {
		panic(err)
	}

}

func (c *ProjectsHandler) ListPages(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf(
		"list project (%s) pages",
		r.Context().Value(projectIDContextKey).(string),
	)))

	if err != nil {
		panic(err)
	}

}

func NewProjectsHandler(l zerolog.Logger, c ProjectsController, jU jsonUnmarshaller) *ProjectsHandler {
	return &ProjectsHandler{
		logger:           l,
		controller:       c,
		jsonUnmarshaller: jU,
	}
}
