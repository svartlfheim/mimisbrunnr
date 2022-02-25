package result_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

func Test_Status_Equals(t *testing.T) {
	tests := []struct{
		name string
		status result.Status
		shouldNotEqual []result.Status
		shouldEqual result.Status
	}{
		{
			name: "created constant equals itself and no others",
			status: result.Created,
			shouldNotEqual: []result.Status{
				result.InternalError,
				result.Invalid,
			},
			shouldEqual: result.Created,
		},
		{
			name: "typecast string created equals constant and no others",
			status: result.Status("created"),
			shouldNotEqual: []result.Status{
				result.InternalError,
				result.Invalid,
			},
			shouldEqual: result.Created,
		},

		{
			name: "internalError constant equals itself and no others",
			status: result.InternalError,
			shouldNotEqual: []result.Status{
				result.Created,
				result.Invalid,
			},
			shouldEqual: result.InternalError,
		},
		{
			name: "typecast string internalError equals constant and no others",
			status: result.Status("internal_error"),
			shouldNotEqual: []result.Status{
				result.Created,
				result.Invalid,
			},
			shouldEqual: result.InternalError,
		},

		{
			name: "invalid constant equals itself and no others",
			status: result.Invalid,
			shouldNotEqual: []result.Status{
				result.Created,
				result.InternalError,
			},
			shouldEqual: result.Invalid,
		},
		{
			name: "typecast string invalid equals constant and no others",
			status: result.Status("invalid"),
			shouldNotEqual: []result.Status{
				result.Created,
				result.InternalError,
			},
			shouldEqual: result.Invalid,
		},

		{
			name: "unknown status is only considered equal to itself",
			status: result.Status("invalid"),
			shouldNotEqual: []result.Status{
				result.Status("garbage"),
				result.Status("other_garbage"),
				result.Status("not_a_status"),
			},
			shouldEqual: result.Invalid,
		},
	}

	for _, test := range(tests) {
		t.Run(test.name, func (tt *testing.T) {
			assert.True(tt, test.status.Equals(test.shouldEqual), "not equal to itself")

			for _, other := range(test.shouldNotEqual) {
				assert.False(tt, test.status.Equals(other), "status '%s' was equal to '%s'", string(test.status), string(other))
			}
		})
	}
}

func Test_Status_ToHTTP(t *testing.T) {
	tests := []struct{
		name string
		status result.Status
		expected int
	}{
		{
			name: "Created status returns 201 created HTTP code",
			status: result.Created,
			expected: 201,
		},
		{
			name: "Internal error status returns 500 created HTTP code",
			status: result.InternalError,
			expected: 500,
		},
		{
			name: "Invalid status returns 422 created HTTP code",
			status: result.Invalid,
			expected: 422,
		},
		{
			name: "Unhandled status returns 501 HTTP code",
			status: result.Status("garbage"),
			expected: 501,
		},
	}

	for _, test := range(tests) {
		t.Run(test.name, func (tt *testing.T) {
			assert.Equal(tt, test.status.ToHTTP(), test.expected)
		})
	}
}