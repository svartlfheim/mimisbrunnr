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

func Test_SCMAccessToken_Getters(t *testing.T) {
	id := uuid.New()
	tCreate := time.Now()
	tUpdate := time.Now().Add(1 * time.Hour)

	at := models.NewSCMAccessToken(id, "myname", "mytoken", true, tCreate, tUpdate)

	assert.Equal(t, id, at.GetID())
	assert.Equal(t, "myname", at.GetName())
	assert.Equal(t, "mytoken", at.GetToken())
	assert.True(t, at.IsActive())
	assert.Equal(t, tCreate, at.GetCreationTime())
	assert.Equal(t, tUpdate, at.GetLastUpdatedTime())
}

func Test_SCMIntegration_SimpleGetters(t *testing.T) {
	id := uuid.New()
	tCreate := time.Now()
	tUpdate := time.Now().Add(1 * time.Hour)

	// Timestamps don't matter here, this is simpler
	tokens := []*models.SCMAccessToken{
		models.NewSCMAccessToken(uuid.New(), "token1", "mytoken1", false, tCreate, tUpdate),
		models.NewSCMAccessToken(uuid.New(), "token2", "mytoken2", false, tCreate, tUpdate),
		models.NewSCMAccessToken(uuid.New(), "token3", "mytoken3", false, tCreate, tUpdate),
	}

	si := models.NewSCMIntegration(id, "myname", models.GithubType, "http://example.com", tokens, tCreate, tUpdate)

	assert.Equal(t, id, si.GetID())
	assert.Equal(t, "myname", si.GetName())
	assert.Equal(t, models.GithubType, si.GetType())
	assert.Equal(t, "http://example.com", si.GetEndpoint())
	assert.Equal(t, tokens, si.GetAccessTokens())
	assert.Equal(t, tCreate, si.GetCreationTime())
	assert.Equal(t, tUpdate, si.GetLastUpdatedTime())
}

func Test_SCMIntegration_GetsActiveAccessToken(t *testing.T) {
	id := uuid.New()
	tCreate := time.Now()
	tUpdate := time.Now().Add(1 * time.Hour)

	// Timestamps don't matter here, this is simpler
	tokens := []*models.SCMAccessToken{
		models.NewSCMAccessToken(uuid.New(), "token1", "mytoken1", true, tCreate, tUpdate),
		models.NewSCMAccessToken(uuid.New(), "token2", "mytoken2", false, tCreate, tUpdate),
		models.NewSCMAccessToken(uuid.New(), "token3", "mytoken3", false, tCreate, tUpdate),
	}
	si := models.NewSCMIntegration(id, "myname", models.GithubType, "http://example.com", tokens, tCreate, tUpdate)

	assert.Equal(t, tokens, si.GetAccessTokens())
	assert.Equal(t, tokens[0], si.GetActiveAccessToken())
}

func Test_SCMIntegration_GetsActiveAccessToken_none_set(t *testing.T) {
	id := uuid.New()
	tCreate := time.Now()
	tUpdate := time.Now().Add(1 * time.Hour)

	// Timestamps don't matter here, this is simpler
	tokens := []*models.SCMAccessToken{
		models.NewSCMAccessToken(uuid.New(), "token1", "mytoken1", false, tCreate, tUpdate),
		models.NewSCMAccessToken(uuid.New(), "token2", "mytoken2", false, tCreate, tUpdate),
		models.NewSCMAccessToken(uuid.New(), "token3", "mytoken3", false, tCreate, tUpdate),
	}

	si := models.NewSCMIntegration(id, "myname", models.GithubType, "http://example.com", tokens, tCreate, tUpdate)

	assert.Equal(t, tokens, si.GetAccessTokens())
	assert.Nil(t, si.GetActiveAccessToken())
}