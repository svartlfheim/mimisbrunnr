package rdb

import (
	"fmt"
	"strings"
)

type ErrConnectionGoneAway struct {
	Retries int
}

func (e ErrConnectionGoneAway) Error() string {
	return fmt.Sprintf("rdb connection has gone away and could not reconnect after %d retries", e.Retries)
}

type ErrUnsupportedDriver struct {
	Driver string
}

func (e ErrUnsupportedDriver) Error() string {
	return fmt.Sprintf("driver '%s' is not supported", e.Driver)
}

type ErrConfigurationMissingFields struct {
	Fields []string
}

func (e ErrConfigurationMissingFields) Error() string {
	return fmt.Sprintf("required fields missing in database configuration: %s", strings.Join(e.Fields, ", "))
}
