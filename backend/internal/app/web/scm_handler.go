package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/scm/v1"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
)

type scmResponseBuilder interface {
	FromCommandResult(apiVersion int, res commandresult.Result) WritableResponse
	FromUnmarshalError(err error) WritableResponse
}

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
	rb               scmResponseBuilder
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

func (h *SCMHandler) RouteGroup() string {
	return "scm-integrations"
}

func (h *SCMHandler) List(w http.ResponseWriter, r *http.Request) {
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

	resp := h.rb.FromCommandResult(
		r.Context().Value(ApiVersionContextKey).(int),
		h.controller.ListV1(dto),
	)

	resp.Write(w)
}

func (h *SCMHandler) Create(w http.ResponseWriter, r *http.Request) {
	dto := v1.AddIntegrationDTO{}
	var resp WritableResponse

	if err := h.jsonUnmarshaller.Unmarshal(r, &dto); err != nil {
		resp = h.rb.FromUnmarshalError(err)
	} else {
		resp = h.rb.FromCommandResult(
			r.Context().Value(ApiVersionContextKey).(int),
			h.controller.AddV1(dto),
		)
	}

	resp.Write(w)
}

func (h *SCMHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(scmIntegrationIDContextKey).(string)

	resp := h.rb.FromCommandResult(
		r.Context().Value(ApiVersionContextKey).(int),
		h.controller.GetV1(id),
	)

	resp.Write(w)
}

func (h *SCMHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(scmIntegrationIDContextKey).(string)
	dto := v1.UpdateIntegrationDTO{}
	var resp WritableResponse

	if err := h.jsonUnmarshaller.Unmarshal(r, &dto); err != nil {
		resp = h.rb.FromUnmarshalError(err)
	} else {
		resp = h.rb.FromCommandResult(
			r.Context().Value(ApiVersionContextKey).(int),
			h.controller.UpdateV1(id, dto),
		)
	}

	resp.Write(w)
}

func (h *SCMHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(scmIntegrationIDContextKey).(string)
	resp := h.rb.FromCommandResult(
		r.Context().Value(ApiVersionContextKey).(int),
		h.controller.DeleteV1(id),
	)

	resp.Write(w)
}

func NewSCMIntegrationsHandler(l zerolog.Logger, m SCMIntegrationsController, jU jsonUnmarshaller, rb scmResponseBuilder) *SCMHandler {
	return &SCMHandler{
		logger:           l,
		controller:       m,
		jsonUnmarshaller: jU,
		rb:               rb,
	}
}
