package scm

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

type AccessToken struct {
	ID uuid.UUID
	Name string
	Token string
	Active bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *AccessToken) GetID() uuid.UUID {
	return c.ID
}

func (c *AccessToken) GetName() string {
	return c.Name
}

func (c *AccessToken) GetToken() string {
	return c.Token
}

func (c *AccessToken) IsActive() bool {
	return c.Active
}

func (c *AccessToken) GetCreationTime() time.Time {
	return c.CreatedAt
}

func (c *AccessToken) GetLastUpdatedTime() time.Time {
	return c.UpdatedAt
}

type SCMIntegration struct {
	ID uuid.UUID
	Name string
	Type SCMIntegrationType
	Endpoint string
	AccessTokens []*AccessToken
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (gh *SCMIntegration) GetID() uuid.UUID {
	return gh.ID
}

func (gh *SCMIntegration) GetName() string {
	return gh.Name
}

func (gh *SCMIntegration) GetType() SCMIntegrationType {
	return gh.Type
}

func (gh *SCMIntegration) GetEndpoint() string {
	return gh.Endpoint
}

func (gh *SCMIntegration) GetCredentials() []*AccessToken {
	return gh.AccessTokens
}

func (gh *SCMIntegration) GetActiveCredentials() *AccessToken {
	for _, c := range(gh.AccessTokens) {
		if c.IsActive() {
			return c
		}
	}

	return nil
}

func (gh *SCMIntegration) GetCreationTime() time.Time {
	return gh.CreatedAt
}

func (gh *SCMIntegration) GetLastUpdatedTime() time.Time {
	return gh.UpdatedAt
}