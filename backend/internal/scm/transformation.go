package scm

import (
	"time"

	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

type scmIntegrationTransformerV1 interface {
	SCMIntegrationV1(m *models.SCMIntegration) *scmIntegrationV1
}

type scmIntegrationV1 struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Endpoint  string `json:"endpoint"`
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Transformer struct {}

func (*Transformer) SCMIntegrationV1(m *models.SCMIntegration) *scmIntegrationV1 {
	return &scmIntegrationV1{
		ID: m.GetID().String(),
		Name: m.GetName(),
		Type: string(m.GetType()),
		Endpoint: m.GetEndpoint(),
		Token: m.GetToken(),
		CreatedAt: m.GetCreationTime().UTC().Format(time.RFC3339),
		UpdatedAt: m.GetLastUpdatedTime().UTC().Format(time.RFC3339),
	}
}

func NewTransformer() *Transformer {
	return &Transformer{}
}