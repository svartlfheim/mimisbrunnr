package api

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

func Test_NewTransformer(t *testing.T) {
	trans := NewTransformer()

	assert.NotNil(t, trans)
}

func Test_Transformer_v1_unsupported_type(t *testing.T) {
	trans := &Transformer{}
	out, err := trans.Transform(1, time.Now())

	require.Nil(t, out)
	require.NotNil(t, err)
	assert.IsType(t, ErrUnsupportedResourceType{}, err)
}

func Test_Transformer_v1_unsupported_api_version(t *testing.T) {
	trans := &Transformer{}
	// Doubt we'll ever get this high...
	out, err := trans.Transform(928346, time.Now())

	require.Nil(t, out)
	require.NotNil(t, err)
	assert.IsType(t, ErrUnsupportedApiVersion{}, err)
}

func Test_Transformer_v1_scmIntegration(t *testing.T) {
	id1 := uuid.New()
	createdAt1 := time.Now()
	updatedAt1 := createdAt1.Add(1 * time.Hour)

	in := models.NewSCMIntegration(
		id1,
		"myintegration",
		models.GithubType,
		"myendpoint.com",
		"mytoken",
		createdAt1,
		updatedAt1,
	)
	expected := &scmIntegrationV1{
		ID:        id1.String(),
		Name:      "myintegration",
		Type:      string(models.GithubType),
		Endpoint:  "myendpoint.com",
		Token:     "mytoken",
		CreatedAt: createdAt1.UTC().Format(time.RFC3339),
		UpdatedAt: updatedAt1.UTC().Format(time.RFC3339),
	}

	trans := &Transformer{}
	out, err := trans.Transform(1, in)

	assert.Equal(t, expected, out)
	assert.Nil(t, err)
}

func Test_Transformer_v1_scmIntegration_list(t *testing.T) {
	id1 := uuid.New()
	createdAt1 := time.Now()
	updatedAt1 := createdAt1.Add(1 * time.Hour)

	id2 := uuid.New()
	createdAt2 := time.Now().Add(1 * time.Minute)
	updatedAt2 := createdAt2.Add(1 * time.Hour)

	in := []*models.SCMIntegration{
		models.NewSCMIntegration(
			id1,
			"myintegration",
			models.GithubType,
			"myendpoint.com",
			"mytoken",
			createdAt1,
			updatedAt1,
		),
		models.NewSCMIntegration(
			id2,
			"myintegration2",
			models.GitlabType,
			"myendpoint2.com",
			"mytoken2",
			createdAt2,
			updatedAt2,
		),
	}
	expected := []*scmIntegrationV1{
		{
			ID:        id1.String(),
			Name:      "myintegration",
			Type:      string(models.GithubType),
			Endpoint:  "myendpoint.com",
			Token:     "mytoken",
			CreatedAt: createdAt1.UTC().Format(time.RFC3339),
			UpdatedAt: updatedAt1.UTC().Format(time.RFC3339),
		},
		{
			ID:        id2.String(),
			Name:      "myintegration2",
			Type:      string(models.GitlabType),
			Endpoint:  "myendpoint2.com",
			Token:     "mytoken2",
			CreatedAt: createdAt2.UTC().Format(time.RFC3339),
			UpdatedAt: updatedAt2.UTC().Format(time.RFC3339),
		},
	}

	trans := &Transformer{}
	out, err := trans.Transform(1, in)

	assert.Equal(t, expected, out)
	assert.Nil(t, err)
}

func Test_Transformer_v1_project(t *testing.T) {
	id1 := uuid.New()
	createdAt1 := time.Now()
	updatedAt1 := createdAt1.Add(1 * time.Hour)

	id2 := uuid.New()
	createdAt2 := time.Now().Add(1 * time.Minute)
	updatedAt2 := createdAt2.Add(1 * time.Hour)

	in := models.NewProject(
		id1,
		"myproject",
		"myorg/myrepo",
		models.NewSCMIntegration(
			id2,
			"myintegration",
			models.GithubType,
			"myendpoint.com",
			"mytoken",
			createdAt2,
			updatedAt2,
		),
		createdAt1,
		updatedAt1,
	)
	expected := &projectV1{
		ID:   id1.String(),
		Name: "myproject",
		Path: "myorg/myrepo",
		SCMIntegration: &scmIntegrationV1{
			ID:        id2.String(),
			Name:      "myintegration",
			Type:      string(models.GithubType),
			Endpoint:  "myendpoint.com",
			Token:     "mytoken",
			CreatedAt: createdAt2.UTC().Format(time.RFC3339),
			UpdatedAt: updatedAt2.UTC().Format(time.RFC3339),
		},
		CreatedAt: createdAt1.UTC().Format(time.RFC3339),
		UpdatedAt: updatedAt1.UTC().Format(time.RFC3339),
	}

	trans := &Transformer{}
	out, err := trans.Transform(1, in)

	assert.Equal(t, expected, out)
	assert.Nil(t, err)
}

func Test_Transformer_v1_project_list(t *testing.T) {
	id1 := uuid.New()
	createdAt1 := time.Now()
	updatedAt1 := createdAt1.Add(1 * time.Hour)

	id2 := uuid.New()
	createdAt2 := time.Now().Add(5 * time.Minute)
	updatedAt2 := createdAt2.Add(1 * time.Hour)

	id3 := uuid.New()
	createdAt3 := time.Now().Add(10 * time.Minute)
	updatedAt3 := createdAt2.Add(1 * time.Hour)

	in := []*models.Project{
		models.NewProject(
			id1,
			"myproject",
			"myorg/myrepo",
			models.NewSCMIntegration(
				id2,
				"myintegration",
				models.GithubType,
				"myendpoint.com",
				"mytoken",
				createdAt2,
				updatedAt2,
			),
			createdAt1,
			updatedAt1,
		),
		models.NewProject(
			id3,
			"myotherproject",
			"myorg/myotherrepo",
			models.NewSCMIntegration(
				id2,
				"myintegration",
				models.GithubType,
				"myendpoint.com",
				"mytoken",
				createdAt2,
				updatedAt2,
			),
			createdAt3,
			updatedAt3,
		),
	}
	expected := []*projectV1{
		{
			ID:   id1.String(),
			Name: "myproject",
			Path: "myorg/myrepo",
			SCMIntegration: &scmIntegrationV1{
				ID:        id2.String(),
				Name:      "myintegration",
				Type:      string(models.GithubType),
				Endpoint:  "myendpoint.com",
				Token:     "mytoken",
				CreatedAt: createdAt2.UTC().Format(time.RFC3339),
				UpdatedAt: updatedAt2.UTC().Format(time.RFC3339),
			},
			CreatedAt: createdAt1.UTC().Format(time.RFC3339),
			UpdatedAt: updatedAt1.UTC().Format(time.RFC3339),
		},
		{
			ID:   id3.String(),
			Name: "myotherproject",
			Path: "myorg/myotherrepo",
			SCMIntegration: &scmIntegrationV1{
				ID:        id2.String(),
				Name:      "myintegration",
				Type:      string(models.GithubType),
				Endpoint:  "myendpoint.com",
				Token:     "mytoken",
				CreatedAt: createdAt2.UTC().Format(time.RFC3339),
				UpdatedAt: updatedAt2.UTC().Format(time.RFC3339),
			},
			CreatedAt: createdAt3.UTC().Format(time.RFC3339),
			UpdatedAt: updatedAt3.UTC().Format(time.RFC3339),
		},
	}

	trans := &Transformer{}
	out, err := trans.Transform(1, in)

	assert.Equal(t, expected, out)
	assert.Nil(t, err)
}
