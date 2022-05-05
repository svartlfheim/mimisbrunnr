package models

import (
	"time"

	"github.com/google/uuid"
)

type SCMIntegrationType string

const GithubType SCMIntegrationType = "github"
const GitlabType SCMIntegrationType = "gitlab"

var availableSCMIntegrationTypes []SCMIntegrationType = []SCMIntegrationType{
	GithubType,
	GitlabType,
}

func AvailableSCMIntegrationTypes() []SCMIntegrationType {
	return availableSCMIntegrationTypes
}

type SCMIntegration struct {
	ID        uuid.UUID
	Name      string
	Type      SCMIntegrationType
	Endpoint  string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (si *SCMIntegration) GetID() uuid.UUID {
	return si.ID
}

func (si *SCMIntegration) GetName() string {
	return si.Name
}

func (si *SCMIntegration) GetType() SCMIntegrationType {
	return si.Type
}

func (si *SCMIntegration) GetEndpoint() string {
	return si.Endpoint
}

func (si *SCMIntegration) GetToken() string {
	return si.Token
}

func (si *SCMIntegration) GetCreationTime() time.Time {
	return si.CreatedAt
}

func (si *SCMIntegration) GetLastUpdatedTime() time.Time {
	return si.UpdatedAt
}

func NewSCMIntegration(id uuid.UUID, name string, t SCMIntegrationType, endpoint string, token string, createdAt time.Time, updatedAt time.Time) *SCMIntegration {
	return &SCMIntegration{
		ID:        id,
		Name:      name,
		Type:      t,
		Endpoint:  endpoint,
		Token:     token,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
