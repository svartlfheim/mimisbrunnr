package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID             uuid.UUID
	Name           string
	Path           string
	Pages          []*ProjectPage
	SCMIntegration *SCMIntegration
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (p *Project) GetID() uuid.UUID {
	return p.ID
}

func (p *Project) GetName() string {
	return p.Name
}

func (p *Project) GetPath() string {
	return p.Path
}

func (p *Project) GetPages() []*ProjectPage {
	return p.Pages
}

func (p *Project) GetSCMIntegration() *SCMIntegration {
	return p.SCMIntegration
}

func (p *Project) GetCreationTime() time.Time {
	return p.CreatedAt
}

func (p *Project) GetLastUpdatedTime() time.Time {
	return p.UpdatedAt
}

func NewProject(id uuid.UUID, name string, path string, SCMIntegration *SCMIntegration, createdAt time.Time, updatedAt time.Time) *Project {
	return &Project{
		ID:             id,
		Name:           name,
		Path:           path,
		Pages:          []*ProjectPage{},
		SCMIntegration: SCMIntegration,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

type ProjectPage struct {
	ID        uuid.UUID
	Title     string
	Path      string
	Parent    *ProjectPage
	Project   *Project
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *ProjectPage) GetID() uuid.UUID {
	return p.ID
}

func (p *ProjectPage) GetTitle() string {
	return p.Title
}

func (p *ProjectPage) GetPath() string {
	return p.Path
}

func (p *ProjectPage) GetParent() *ProjectPage {
	return p.Parent
}

func (p *ProjectPage) GetProject() *Project {
	return p.Project
}

func (p *ProjectPage) GetCreationTime() time.Time {
	return p.CreatedAt
}

func (p *ProjectPage) GetLastUpdatedTime() time.Time {
	return p.UpdatedAt
}
