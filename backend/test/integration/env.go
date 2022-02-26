package integration

import (
	"fmt"
	"os"
)

const EnvPrefix string = "ITEST_MIMISBRUNNR_"

func Getenv(k string) (string, bool) {
	k = fmt.Sprintf("%s%s", EnvPrefix, k)
	return os.LookupEnv(k)
}
