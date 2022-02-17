package githosts

import (
	"time"

	"github.com/google/uuid"
)

type GitHostType string

const GithubType GitHostType = "github"
const GitlabType GitHostType = "gitlab"
var availableGitHostTypes []GitHostType = []GitHostType{
	GithubType,
	GitlabType,
}

func AvailableGitHostTypes() []GitHostType {
	return availableGitHostTypes
}

type Credentials struct {
	ID uuid.UUID
	PersonalAccessToken string
	Active bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Credentials) GetID() uuid.UUID {
	return c.ID
}

func (c *Credentials) GetToken() string {
	return c.PersonalAccessToken
}

func (c *Credentials) IsActive() bool {
	return c.Active
}

func (c *Credentials) GetCreationTime() time.Time {
	return c.CreatedAt
}

func (c *Credentials) GetLastUpdatedTime() time.Time {
	return c.UpdatedAt
}

type GitHost struct {
	ID uuid.UUID
	Name string
	Type GitHostType
	Endpoint string
	Credentials []*Credentials
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (gh *GitHost) GetID() uuid.UUID {
	return gh.ID
}

func (gh *GitHost) GetName() string {
	return gh.Name
}

func (gh *GitHost) GetType() GitHostType {
	return gh.Type
}

func (gh *GitHost) GetEndpoint() string {
	return gh.Endpoint
}

func (gh *GitHost) GetCredentials() []*Credentials {
	return gh.Credentials
}

func (gh *GitHost) GetActiveCredentials() *Credentials {
	for _, c := range(gh.Credentials) {
		if c.IsActive() {
			return c
		}
	}

	return nil
}

func (gh *GitHost) GetCreationTime() time.Time {
	return gh.CreatedAt
}

func (gh *GitHost) GetLastUpdatedTime() time.Time {
	return gh.UpdatedAt
}