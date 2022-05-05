package testenv

import "os"

type ResetFunc func() error

func Override(k string, v string) ResetFunc {
	val, found := os.LookupEnv(k)

	os.Setenv(k, v)

	if !found {
		return ResetFunc(func() error {
			return os.Unsetenv(k)
		})
	}

	return ResetFunc(func() error {
		return os.Setenv(k, val)
	})
}

func Clear(keys ...string) error {
	for _, k := range keys {
		if err := os.Unsetenv(k); err != nil {
			return err
		}
	}

	return nil
}
