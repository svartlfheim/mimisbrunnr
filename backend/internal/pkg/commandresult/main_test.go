package commandresult_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
)

func Test_Status_Equals(t *testing.T) {
	tests := []struct {
		name           string
		status         commandresult.Status
		shouldNotEqual []commandresult.Status
		shouldEqual    commandresult.Status
	}{
		{
			name:   "created constant equals itself and no others",
			status: commandresult.Created,
			shouldNotEqual: []commandresult.Status{
				commandresult.InternalError,
				commandresult.Invalid,
			},
			shouldEqual: commandresult.Created,
		},
		{
			name:   "typecast string created equals constant and no others",
			status: commandresult.Status("created"),
			shouldNotEqual: []commandresult.Status{
				commandresult.InternalError,
				commandresult.Invalid,
			},
			shouldEqual: commandresult.Created,
		},

		{
			name:   "internalError constant equals itself and no others",
			status: commandresult.InternalError,
			shouldNotEqual: []commandresult.Status{
				commandresult.Created,
				commandresult.Invalid,
			},
			shouldEqual: commandresult.InternalError,
		},
		{
			name:   "typecast string internalError equals constant and no others",
			status: commandresult.Status("internal_error"),
			shouldNotEqual: []commandresult.Status{
				commandresult.Created,
				commandresult.Invalid,
			},
			shouldEqual: commandresult.InternalError,
		},

		{
			name:   "invalid constant equals itself and no others",
			status: commandresult.Invalid,
			shouldNotEqual: []commandresult.Status{
				commandresult.Created,
				commandresult.InternalError,
			},
			shouldEqual: commandresult.Invalid,
		},
		{
			name:   "typecast string invalid equals constant and no others",
			status: commandresult.Status("invalid"),
			shouldNotEqual: []commandresult.Status{
				commandresult.Created,
				commandresult.InternalError,
			},
			shouldEqual: commandresult.Invalid,
		},

		{
			name:   "unknown status is only considered equal to itself",
			status: commandresult.Status("invalid"),
			shouldNotEqual: []commandresult.Status{
				commandresult.Status("garbage"),
				commandresult.Status("other_garbage"),
				commandresult.Status("not_a_status"),
			},
			shouldEqual: commandresult.Invalid,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert.True(tt, test.status.Equals(test.shouldEqual), "not equal to itself")

			for _, other := range test.shouldNotEqual {
				assert.False(tt, test.status.Equals(other), "status '%s' was equal to '%s'", string(test.status), string(other))
			}
		})
	}
}

func Test_Status_ToHTTP(t *testing.T) {
	tests := []struct {
		name     string
		status   commandresult.Status
		expected int
	}{
		{
			name:     "Created status returns 201 created HTTP code",
			status:   commandresult.Created,
			expected: 201,
		},
		{
			name:     "Internal error status returns 500 created HTTP code",
			status:   commandresult.InternalError,
			expected: 500,
		},
		{
			name:     "Invalid status returns 422 created HTTP code",
			status:   commandresult.Invalid,
			expected: 422,
		},
		{
			name:     "Unhandled status returns 501 HTTP code",
			status:   commandresult.Status("garbage"),
			expected: 501,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert.Equal(tt, test.status.ToHTTP(), test.expected)
		})
	}
}
