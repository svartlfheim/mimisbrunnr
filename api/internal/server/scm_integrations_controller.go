package server

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/internal/scm"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type jsonUnmarshaller interface {
	Unmarshal(contents io.Reader, into interface{}) error
}

type SCMIntegrationsManager interface {
	Add(dto scm.AddSCMIntegrationDTO) (result.Result)
}

type SCMIntegrationsController struct {
	logger zerolog.Logger
	manager SCMIntegrationsManager
	jsonUnmarshaller jsonUnmarshaller
}

func SCMIntegrationContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SCMIntegrationID := chi.URLParam(r, "SCMIntegrationID")

		ctx := context.WithValue(r.Context(), "SCMIntegrationID", SCMIntegrationID)
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
	w.Write([]byte("list scm integrations"))
}

func (c *SCMIntegrationsController) Search(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("search scm integrations"))
}

func (c *SCMIntegrationsController) Create(w http.ResponseWriter, r *http.Request) {
	dto := scm.AddSCMIntegrationDTO{}
	

	if err := c.jsonUnmarshaller.Unmarshal(r.Body, &dto); handleError(w, c.logger, err) {
		return
	}
	res := c.manager.Add(dto)
	msgs := []string{}
	
	for _, err := range(res.Errors()) {
		msgs = append(msgs, err.Error())
	}

	w.WriteHeader(res.Status().ToHTTP())
	w.Write([]byte(strings.Join(msgs, "\n")))
}

func (c *SCMIntegrationsController) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get scm integration: " + r.Context().Value("SCMIntegrationID").(string)))
}

func (c *SCMIntegrationsController) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update scm integration: " + r.Context().Value("SCMIntegrationID").(string)))
}

func (c *SCMIntegrationsController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete scm integration: " + r.Context().Value("SCMIntegrationID").(string)))
}

func NewSCMIntegrationsController(l zerolog.Logger, m SCMIntegrationsManager, jU jsonUnmarshaller) *SCMIntegrationsController {
	return &SCMIntegrationsController{
		logger: l,
		manager: m,
		jsonUnmarshaller: jU,
	}
}