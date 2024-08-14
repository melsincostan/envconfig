package envconfig

import (
	"fmt"
	"os"
	"strings"
)

func getString(env_name string, required bool) (string, error) {
	raw := os.Getenv(env_name)
	val := strings.TrimSpace(raw)
	if required && len(val) < 1 {
		return "", fmt.Errorf("could not find %s even though it is required", env_name)
	}
	return val, nil
}
