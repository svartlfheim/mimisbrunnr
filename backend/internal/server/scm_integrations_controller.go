package server

import (
	"context"
	"net/http"

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
	r.Get("/search", c.Search)

	r.Post("/", c.Create)

	r.Route("/{SCMIntegrationID}", func(r chi.Router) {
		r.Use(SCMIntegrationContext)
		r.Get("/", c.Get)
		r.Put("/", c.Update)
		r.Delete("/", c.Delete)
	})

	return r
}

func (c *SCMIntegrationsController) RouteGroup() string {
	return "scm-integrations"
}

func (c *SCMIntegrationsController) List(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("list scm integrations"))

	if err != nil {
		panic(err)
	}
}

func (c *SCMIntegrationsController) Search(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("search scm integrations"))

	if err != nil {
		panic(err)
	}
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
	_, err := w.Write([]byte("get scm integration: " + r.Context().Value(scmIntegrationIDContextKey).(string)))

	if err != nil {
		panic(err)
	}
}

func (c *SCMIntegrationsController) Update(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("update scm integration: " + r.Context().Value(scmIntegrationIDContextKey).(string)))

	if err != nil {
		panic(err)
	}
}

func (c *SCMIntegrationsController) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("delete scm integration: " + r.Context().Value(scmIntegrationIDContextKey).(string)))

	if err != nil {
		panic(err)
	}
}

func NewSCMIntegrationsController(l zerolog.Logger, m SCMIntegrationsManager, jU jsonUnmarshaller) *SCMIntegrationsController {
	return &SCMIntegrationsController{
		logger:           l,
		manager:          m,
		jsonUnmarshaller: jU,
	}
}
