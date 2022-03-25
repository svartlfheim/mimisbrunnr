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

type projectsResponseBuilder interface {
	FromCommandResult(apiVersion int, res commandresult.Result) WritableResponse
	FromUnmarshalError(err error) WritableResponse
}

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
	rb projectsResponseBuilder
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

func (h *ProjectsHandler) RouteGroup() string {
	return "projects"
}

func (h *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
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

	resp := h.rb.FromCommandResult(
		r.Context().Value(apiVersionContextKey).(int),
		h.controller.ListV1(dto),
	)

	resp.Write(w)
}

func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	dto := v1.AddProjectDTO{}
	var resp WritableResponse

	if err := h.jsonUnmarshaller.Unmarshal(r, &dto); err != nil {
		resp = h.rb.FromUnmarshalError(err)
	} else {
		resp = h.rb.FromCommandResult(
			r.Context().Value(apiVersionContextKey).(int),
			h.controller.AddV1(dto),
		)
	}

	resp.Write(w)
}

func (h *ProjectsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(projectIDContextKey).(string)

	resp := h.rb.FromCommandResult(
		r.Context().Value(apiVersionContextKey).(int),
		h.controller.GetV1(id),
	)

	resp.Write(w)
}

func (h *ProjectsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(projectIDContextKey).(string)
	dto := v1.UpdateProjectDTO{}
	var resp WritableResponse

	if err := h.jsonUnmarshaller.Unmarshal(r, &dto); err != nil {
		resp = h.rb.FromUnmarshalError(err)
	} else {
		resp = h.rb.FromCommandResult(
			r.Context().Value(apiVersionContextKey).(int),
			h.controller.UpdateV1(id, dto),
		)
	}

	resp.Write(w)
}

func (h *ProjectsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(projectIDContextKey).(string)
	resp := h.rb.FromCommandResult(
		r.Context().Value(apiVersionContextKey).(int),
		h.controller.DeleteV1(id),
	)

	resp.Write(w)
}

func (h *ProjectsHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf(
		"get project (%s) page: %s",
		r.Context().Value(projectIDContextKey).(string),
		r.Context().Value(pageIDContextKey).(string),
	)))

	if err != nil {
		panic(err)
	}
}

func (h *ProjectsHandler) SearchPages(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf(
		"search project (%s) pages",
		r.Context().Value(projectIDContextKey).(string),
	)))

	if err != nil {
		panic(err)
	}

}

func (h *ProjectsHandler) ListPages(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf(
		"list project (%s) pages",
		r.Context().Value(projectIDContextKey).(string),
	)))

	if err != nil {
		panic(err)
	}

}

func NewProjectsHandler(l zerolog.Logger, c ProjectsController, jU jsonUnmarshaller, rb projectsResponseBuilder) *ProjectsHandler {
	return &ProjectsHandler{
		logger:           l,
		controller:       c,
		jsonUnmarshaller: jU,
		rb:            rb,
	}
}
