package env

import (
	"errors"
	"fmt"
	"os"
)

func Check(envs ...string) (err error) {
	for _, e := range envs {
		if _, exists := os.LookupEnv(e); !exists {
			err = errors.Join(
				err,
				fmt.Errorf("env %s not found", e),
			)
		}
	}
	return err
}
