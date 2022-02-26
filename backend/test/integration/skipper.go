package integration

import (
	"os"
	"testing"
)

func SkipIfIntegrationTestsNotConfigured(t *testing.T) {
	if val, found := os.LookupEnv("CI_INTEGRATION_TESTS_ENABLED"); !found || val != "true" {
		t.Skip("Skipping integration test, set CI_INTEGRATION_TESTS_ENABLED=true to run this test")
	}
}
