package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/internal/scm"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type jsonUnmarshaller interface {
	Unmarshal(r *http.Request, into interface{}) error
}

type SCMIntegrationsManager interface {
	AddV1(dto scm.AddSCMIntegrationV1DTO) result.Result
	GetV1(id string) result.Result
	ListV1(dto scm.ListSCMIntegrationsV1DTO) result.Result
	UpdateV1(id string, dto scm.UpdateSCMIntegrationV1DTO) result.Result
	DeleteV1(id string) result.Result
}

type SCMIntegrationsController struct {
	logger           zerolog.Logger
	manager          SCMIntegrationsManager
	jsonUnmarshaller jsonUnmarshaller
}

func SCMIntegrationContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SCMIntegrationID := chi.URLParam(r, "SCMIntegrationID")

		ctx := context.WithValue(r.Context(), scmIntegrationIDContextKey, SCMIntegrationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *SCMIntegrationsController) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", c.List)
	r.Post("/", c.Create)

	r.Route("/{SCMIntegrationID}", func(r chi.Router) {
		r.Use(SCMIntegrationContext)
		r.Get("/", c.Get)
		r.Patch("/", c.Update)
		r.Delete("/", c.Delete)
	})

	return r
}

func (c *SCMIntegrationsController) RouteGroup() string {
	return "scm-integrations"
}

func (c *SCMIntegrationsController) List(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	dto := scm.ListSCMIntegrationsV1DTO{}

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
		c.manager.ListV1(dto),
		c.logger,
		w,
	)
}

func (c *SCMIntegrationsController) Create(w http.ResponseWriter, r *http.Request) {
	dto := scm.AddSCMIntegrationV1DTO{}

	if err := c.jsonUnmarshaller.Unmarshal(r, &dto); handleError(w, c.logger, err) {
		return
	}

	serveResponseForResult(
		c.manager.AddV1(dto),
		c.logger,
		w,
	)
}

func (c *SCMIntegrationsController) Get(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(scmIntegrationIDContextKey).(string)

	serveResponseForResult(
		c.manager.GetV1(id),
		c.logger,
		w,
	)
}

func (c *SCMIntegrationsController) Update(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(scmIntegrationIDContextKey).(string)

	dto := scm.UpdateSCMIntegrationV1DTO{}

	if err := c.jsonUnmarshaller.Unmarshal(r, &dto); handleError(w, c.logger, err) {
		return
	}

	serveResponseForResult(
		c.manager.UpdateV1(id, dto),
		c.logger,
		w,
	)
}

func (c *SCMIntegrationsController) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(scmIntegrationIDContextKey).(string)

	serveResponseForResult(
		c.manager.DeleteV1(id),
		c.logger,
		w,
	)
}

func NewSCMIntegrationsController(l zerolog.Logger, m SCMIntegrationsManager, jU jsonUnmarshaller) *SCMIntegrationsController {
	return &SCMIntegrationsController{
		logger:           l,
		manager:          m,
		jsonUnmarshaller: jU,
	}
}
