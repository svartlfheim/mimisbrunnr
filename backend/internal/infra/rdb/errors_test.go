package rdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/rdb"
)

func Test_ErrConnectionGoneAway(t *testing.T) {
	err := rdb.ErrConnectionGoneAway{
		Retries: 5,
	}

	assert.Equal(
		t,
		"rdb connection has gone away and could not reconnect after 5 retries",
		err.Error(),
	)
}

func Test_ErrUnsupportedDriver(t *testing.T) {
	err := rdb.ErrUnsupportedDriver{
		Driver: "baddriver",
	}

	assert.Equal(
		t,
		"driver 'baddriver' is not supported",
		err.Error(),
	)
}

func Test_ErrConfigurationMissingFields(t *testing.T) {
	err := rdb.ErrConfigurationMissingFields{
		Fields: []string{"field1", "field2", "field3"},
	}

	assert.Equal(
		t,
		"required fields missing in database configuration: field1, field2, field3",
		err.Error(),
	)
}
