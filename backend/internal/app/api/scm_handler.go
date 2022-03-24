package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/scm/v1"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
)

type SCMIntegrationsController interface {
	AddV1(dto v1.AddIntegrationDTO) commandresult.Result
	GetV1(id string) commandresult.Result
	ListV1(dto v1.ListIntegrationsDTO) commandresult.Result
	UpdateV1(id string, dto v1.UpdateIntegrationDTO) commandresult.Result
	DeleteV1(id string) commandresult.Result
}

type SCMHandler struct {
	logger           zerolog.Logger
	controller       SCMIntegrationsController
	jsonUnmarshaller jsonUnmarshaller
}

func SCMIntegrationContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SCMIntegrationID := chi.URLParam(r, "SCMIntegrationID")

		ctx := context.WithValue(r.Context(), scmIntegrationIDContextKey, SCMIntegrationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *SCMHandler) Routes() http.Handler {
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

func (c *SCMHandler) RouteGroup() string {
	return "scm-integrations"
}

func (c *SCMHandler) List(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	dto := v1.ListIntegrationsDTO{}

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

func (c *SCMHandler) Create(w http.ResponseWriter, r *http.Request) {
	dto := v1.AddIntegrationDTO{}

	if err := c.jsonUnmarshaller.Unmarshal(r, &dto); handleError(w, c.logger, err) {
		return
	}

	serveResponseForResult(
		c.controller.AddV1(dto),
		c.logger,
		w,
	)
}

func (c *SCMHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(scmIntegrationIDContextKey).(string)

	serveResponseForResult(
		c.controller.GetV1(id),
		c.logger,
		w,
	)
}

func (c *SCMHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(scmIntegrationIDContextKey).(string)

	dto := v1.UpdateIntegrationDTO{}

	if err := c.jsonUnmarshaller.Unmarshal(r, &dto); handleError(w, c.logger, err) {
		return
	}

	serveResponseForResult(
		c.controller.UpdateV1(id, dto),
		c.logger,
		w,
	)
}

func (c *SCMHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(scmIntegrationIDContextKey).(string)

	serveResponseForResult(
		c.controller.DeleteV1(id),
		c.logger,
		w,
	)
}

func NewSCMIntegrationsHandler(l zerolog.Logger, m SCMIntegrationsController, jU jsonUnmarshaller) *SCMHandler {
	return &SCMHandler{
		logger:           l,
		controller:       m,
		jsonUnmarshaller: jU,
	}
}
