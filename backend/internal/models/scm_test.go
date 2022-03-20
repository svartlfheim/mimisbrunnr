package models_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

func Test_AvailableSCMIntegrationTypes(t *testing.T) {
	assert.Equal(t, []models.SCMIntegrationType{
		models.SCMIntegrationType("github"),
		models.SCMIntegrationType("gitlab"),
	}, models.AvailableSCMIntegrationTypes())
}

func Test_SCMIntegration_SimpleGetters(t *testing.T) {
	id := uuid.New()
	tCreate := time.Now()
	tUpdate := time.Now().Add(1 * time.Hour)

	si := models.NewSCMIntegration(id, "myname", models.GithubType, "http://example.com", "mytoken", tCreate, tUpdate)

	assert.Equal(t, id, si.GetID())
	assert.Equal(t, "myname", si.GetName())
	assert.Equal(t, models.GithubType, si.GetType())
	assert.Equal(t, "http://example.com", si.GetEndpoint())
	assert.Equal(t, "mytoken", si.GetToken())
	assert.Equal(t, tCreate, si.GetCreationTime())
	assert.Equal(t, tUpdate, si.GetLastUpdatedTime())
}
