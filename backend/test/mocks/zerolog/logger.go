package zerologmocks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type Zerologger struct {
	Buffer *bytes.Buffer
	Logger zerolog.Logger
}

type filterOpt func(map[string]interface{})

func IgnoreFieldFilter(f string) filterOpt {
	return func(m map[string]interface{}) {
		delete(m, f)
	}
}

func (zl *Zerologger) ExtractLogsToMap(filters ...filterOpt) (m []map[string]interface{}) {
	lines := strings.Split(zl.Buffer.String(), "\n")

	// trailing newline needs removing
	lines = lines[:len(lines)-1]

	for _, line := range lines {
		into := map[string]interface{}{}
		err := json.Unmarshal([]byte(line), &into)

		if err != nil {
			panic(fmt.Errorf("Failed to unmarshal json: %w", err))
		}

		// Remove the time field if it exists...
		// It tends to be problematic for tests, and very rarely a useful test criterion
		delete(into, "time")

		for _, f := range filters {
			f(into)
		}

		m = append(m, into)
	}

	return
}

func (zl *Zerologger) AssertLogs(t *testing.T, logs []map[string]interface{}, filters ...filterOpt) {
	assert.Equal(t, logs, zl.ExtractLogsToMap(filters...))
}

func NewLogger() *Zerologger {
	b := new(bytes.Buffer)

	return &Zerologger{
		Buffer: b,
		Logger: zerolog.New(b).With().Timestamp().Logger(),
	}
}
