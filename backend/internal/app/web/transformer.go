package web

import (
	"time"

	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

type ProjectV1 struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Path           string           `json:"path"`
	SCMIntegration *SCMIntegrationV1 `json:"scm_integration"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
}

type SCMIntegrationV1 struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Endpoint  string `json:"endpoint"`
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Transformer struct{}

func (t *Transformer) ProjectV1(m *models.Project) *ProjectV1 {
	if m == (*models.Project)(nil) {
		return (*ProjectV1)(nil)
	}

	return &ProjectV1{
		ID:             m.GetID().String(),
		Name:           m.GetName(),
		Path:           string(m.GetPath()),
		SCMIntegration: t.IntegrationV1(m.GetSCMIntegration()),
		CreatedAt:      m.GetCreationTime().UTC().Format(time.RFC3339),
		UpdatedAt:      m.GetLastUpdatedTime().UTC().Format(time.RFC3339),
	}
}

func (t *Transformer) ProjectListV1(list []*models.Project) []*ProjectV1 {
	if list == nil {
		return []*ProjectV1{}
	}
	transformed := []*ProjectV1{}

	for _, i := range list {
		transformed = append(transformed, t.ProjectV1(i))
	}

	return transformed
}

func (*Transformer) IntegrationV1(m *models.SCMIntegration) *SCMIntegrationV1 {
	if m == (*models.SCMIntegration)(nil) {
		return (*SCMIntegrationV1)(nil)
	}

	return &SCMIntegrationV1{
		ID:        m.GetID().String(),
		Name:      m.GetName(),
		Type:      string(m.GetType()),
		Endpoint:  m.GetEndpoint(),
		Token:     m.GetToken(),
		CreatedAt: m.GetCreationTime().UTC().Format(time.RFC3339),
		UpdatedAt: m.GetLastUpdatedTime().UTC().Format(time.RFC3339),
	}
}

func (t *Transformer) IntegrationListV1(list []*models.SCMIntegration) []*SCMIntegrationV1 {
	if list == nil {
		return []*SCMIntegrationV1{}
	}

	transformed := []*SCMIntegrationV1{}

	for _, i := range list {
		transformed = append(transformed, t.IntegrationV1(i))
	}

	return transformed
}

func (t *Transformer) transformV1(item interface{}) (interface{}, error) {
	switch item.(type) {
	case *models.Project:
		return t.ProjectV1(item.(*models.Project)), nil
	case []*models.Project:
		return t.ProjectListV1(item.([]*models.Project)), nil
	case *models.SCMIntegration:
		return t.IntegrationV1(item.(*models.SCMIntegration)), nil
	case []*models.SCMIntegration:
		return t.IntegrationListV1(item.([]*models.SCMIntegration)), nil
	default:
		return nil, ErrUnsupportedResourceType{
			Val: item,
		}
	}
}

func (t *Transformer) Transform(v int, item interface{}) (interface{}, error) {
	switch v {
	case 1:
		return t.transformV1(item)
	default:
		return nil, ErrUnsupportedApiVersion{
			Version: v,
		}
	}
}

func NewTransformer() *Transformer {
	return &Transformer{}
}
