package projects

import (
	"github.com/rs/zerolog"
	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/projects/v1"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
)

type projectsRepository interface {
	v1.RequiredRepository
}

type integrationRepository interface {
	v1.RequiredIntegrationRepository
}

type projectsTransformer interface {
	v1.Transformer
}

type structValidator interface {
	v1.StructValidator
}

type Controller struct {
	logger      zerolog.Logger
	repo        projectsRepository
	iRepo       integrationRepository
	validator   structValidator
	transformer projectsTransformer
}

func (m *Controller) AddV1(dto v1.AddProjectDTO) commandresult.Result {
	return v1.Add(m.repo, m.iRepo, m.validator, m.transformer, dto)
}

func (m *Controller) ListV1(dto v1.ListProjectsDTO) commandresult.Result {
	return v1.List(m.repo, m.validator, m.transformer, dto)
}

func (m *Controller) GetV1(id string) commandresult.Result {
	return v1.Get(m.repo, m.transformer, id)
}

func (m *Controller) UpdateV1(id string, dto v1.UpdateProjectDTO) commandresult.Result {
	return v1.Update(m.repo, m.iRepo, m.validator, m.transformer, id, dto)
}

func (m *Controller) DeleteV1(id string) commandresult.Result {
	return v1.Delete(m.repo, m.transformer, id)
}

func NewController(l zerolog.Logger, repo projectsRepository, iRepo integrationRepository, v structValidator, t projectsTransformer) *Controller {
	return &Controller{
		logger:      l,
		repo:        repo,
		iRepo:       iRepo,
		validator:   v,
		transformer: t,
	}
}
