package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

const SCMIntegrationsTableName string = "scm_integrations"
const ProjectsTableName string = "projects"

type scmIntegration struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Token     string    `db:"token"`
	Endpoint  string    `db:"endpoint"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (pI scmIntegration) ToDomainModel() *models.SCMIntegration {
	return models.NewSCMIntegration(
		uuid.MustParse(pI.ID),
		pI.Name,
		models.SCMIntegrationType(pI.Type),
		pI.Endpoint,
		pI.Token,
		pI.CreatedAt,
		pI.UpdatedAt,
	)
}

func toDBSCMIntegration(gh *models.SCMIntegration) *scmIntegration {
	return &scmIntegration{
		ID:        gh.GetID().String(),
		Name:      gh.GetName(),
		Type:      string(gh.GetType()),
		Token:     gh.GetToken(),
		Endpoint:  gh.GetEndpoint(),
		CreatedAt: gh.GetCreationTime(),
		UpdatedAt: gh.GetLastUpdatedTime(),
	}
}

type project struct {
	ID             string    `db:"id"`
	Name           string    `db:"name"`
	Path           string    `db:"path"`
	SCMIntegration string    `db:"scm_integration_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (p project) ToDomainModel(i *models.SCMIntegration) *models.Project {
	return models.NewProject(
		uuid.MustParse(p.ID),
		p.Name,
		p.Path,
		i,
		p.CreatedAt,
		p.UpdatedAt,
	)
}

func toDBProject(p *models.Project) *project {
	return &project{
		ID:             p.GetID().String(),
		Name:           p.GetName(),
		Path:           p.GetPath(),
		SCMIntegration: p.GetSCMIntegration().ID.String(),
		CreatedAt:      p.GetCreationTime(),
		UpdatedAt:      p.GetLastUpdatedTime(),
	}
}
