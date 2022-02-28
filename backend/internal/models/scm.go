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

type SCMAccessToken struct {
	ID        uuid.UUID
	Name      string
	Token     string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *SCMAccessToken) GetID() uuid.UUID {
	return c.ID
}

func (c *SCMAccessToken) GetName() string {
	return c.Name
}

func (c *SCMAccessToken) GetToken() string {
	return c.Token
}

func (c *SCMAccessToken) IsActive() bool {
	return c.Active
}

func (c *SCMAccessToken) GetCreationTime() time.Time {
	return c.CreatedAt
}

func (c *SCMAccessToken) GetLastUpdatedTime() time.Time {
	return c.UpdatedAt
}

func NewSCMAccessToken(id uuid.UUID, name string, token string, active bool, createdAt time.Time, updatedAt time.Time) *SCMAccessToken {
	return &SCMAccessToken{
		ID:        id,
		Name:      name,
		Token:     token,
		Active:    active,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type SCMIntegration struct {
	ID           uuid.UUID
	Name         string
	Type         SCMIntegrationType
	Endpoint     string
	AccessTokens []*SCMAccessToken
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

func (si *SCMIntegration) GetAccessTokens() []*SCMAccessToken {
	return si.AccessTokens
}

func (si *SCMIntegration) GetActiveAccessToken() *SCMAccessToken {
	for _, c := range si.AccessTokens {
		if c.IsActive() {
			return c
		}
	}

	return nil
}

func (si *SCMIntegration) GetCreationTime() time.Time {
	return si.CreatedAt
}

func (si *SCMIntegration) GetLastUpdatedTime() time.Time {
	return si.UpdatedAt
}

func NewSCMIntegration(id uuid.UUID, name string, t SCMIntegrationType, endpoint string, accessTokens []*SCMAccessToken, createdAt time.Time, updatedAt time.Time) *SCMIntegration {
	return &SCMIntegration{
		ID:           id,
		Name:         name,
		Type:         t,
		Endpoint:     endpoint,
		AccessTokens: accessTokens,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}
