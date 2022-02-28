package scm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
	validationmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/validation"
)

func Test_Response(t *testing.T) {
	si := &models.SCMIntegration{}
	resp := &AddSCMIntegrationV1Response{
		created: si,
		errors: []error{
			errors.New("blah"),
		},
		status: result.Created,
		validationErrors: []validation.ValidationError{
			&validationmocks.ValidationError{},
		},
	}

	assert.Implements(t, (*result.Result)(nil), resp)
	assert.Same(t, resp.created, resp.Data())
	// Meta is not yet implemented here
	assert.Equal(t, map[string]interface{}{}, resp.Meta())
	assert.Equal(t, resp.errors, resp.Errors())
	assert.Equal(t, resp.validationErrors, resp.ValidationErrors())
	assert.Equal(t, resp.status, resp.Status())
	assert.False(t, resp.IsListData())
}