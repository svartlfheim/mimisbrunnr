package expectations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type ValidationError struct {
	Path  string
	Rule  string
	Param map[string]string
}

func (v ValidationError) Assert(t *testing.T, e validation.ValidationError) {
	assert.Equal(t, v.Path, e.Path())
	assert.Equal(t, v.Rule, e.Rule())
	assert.Equal(t, v.Param, e.Parameters())
}
