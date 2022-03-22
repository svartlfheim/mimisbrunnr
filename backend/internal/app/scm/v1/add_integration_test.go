package v1

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
	validationmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/pkg/validation"
)

func Test_Response(t *testing.T) {
	si := &TransformedSCMIntegration{}
	resp := &addIntegrationV1Response{
		created: si,
		errors: []error{
			errors.New("blah"),
		},
		status: commandresult.Created,
		validationErrors: []validation.ValidationError{
			&validationmocks.ValidationError{},
		},
	}

	assert.Implements(t, (*commandresult.Result)(nil), resp)
	assert.Same(t, resp.created, resp.Data())
	// Meta is not yet implemented here
	assert.Equal(t, map[string]interface{}{}, resp.Meta())
	assert.Equal(t, resp.errors, resp.Errors())
	assert.Equal(t, resp.validationErrors, resp.ValidationErrors())
	assert.Equal(t, resp.status, resp.Status())
	assert.False(t, resp.IsListData())
}
