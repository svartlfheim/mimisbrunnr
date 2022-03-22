package scm

import (
	"time"

	"github.com/svartlfheim/mimisbrunnr/internal/models"
	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/scm/v1"
)

type Transformer struct{}

func (*Transformer) IntegrationV1(m *models.SCMIntegration) *v1.TransformedSCMIntegration {
	return &v1.TransformedSCMIntegration{
		ID:        m.GetID().String(),
		Name:      m.GetName(),
		Type:      string(m.GetType()),
		Endpoint:  m.GetEndpoint(),
		Token:     m.GetToken(),
		CreatedAt: m.GetCreationTime().UTC().Format(time.RFC3339),
		UpdatedAt: m.GetLastUpdatedTime().UTC().Format(time.RFC3339),
	}
}

func (t *Transformer) IntegrationListV1(list []*models.SCMIntegration) []*v1.TransformedSCMIntegration {
	transformed := []*v1.TransformedSCMIntegration{}

	for _, i := range list {
		transformed = append(transformed, t.IntegrationV1(i))
	}

	return transformed
}

func NewTransformer() *Transformer {
	return &Transformer{}
}
