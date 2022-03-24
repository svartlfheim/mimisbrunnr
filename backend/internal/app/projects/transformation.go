package projects

import (
	"time"

	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/projects/v1"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

type Transformer struct{}

func (*Transformer) ProjectV1(m *models.Project) *v1.TransformedProject {
	return &v1.TransformedProject{
		ID:            m.GetID().String(),
		Name:          m.GetName(),
		Path:          string(m.GetPath()),
		IntegrationID: m.GetSCMIntegration().GetID().String(),
		CreatedAt:     m.GetCreationTime().UTC().Format(time.RFC3339),
		UpdatedAt:     m.GetLastUpdatedTime().UTC().Format(time.RFC3339),
	}
}

func (t *Transformer) ProjectListV1(list []*models.Project) []*v1.TransformedProject {
	transformed := []*v1.TransformedProject{}

	for _, i := range list {
		transformed = append(transformed, t.ProjectV1(i))
	}

	return transformed
}

func NewTransformer() *Transformer {
	return &Transformer{}
}
