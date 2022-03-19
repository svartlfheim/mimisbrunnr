package integration

import (
	"fmt"
	"os"
	"testing"
)

const EnvPrefix string = "ITEST_MIMISBRUNNR_"

func Getenv(k string) (string, bool) {
	k = fmt.Sprintf("%s%s", EnvPrefix, k)
	return os.LookupEnv(k)
}

func GetenvOrFail(t *testing.T, k string) (string) {
	v, found := Getenv(k);

	if !found {
		t.Errorf("expected env var to be present: %s", k)
		t.FailNow()
	}

	return v
}