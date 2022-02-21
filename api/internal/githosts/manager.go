package githosts

import (
	"github.com/rs/zerolog"
)

type Manager struct {
	logger zerolog.Logger
	repo repository
}

func NewManager(l zerolog.Logger, repo repository) *Manager {
	return &Manager{
		logger: l,
		repo: repo,
	}
}