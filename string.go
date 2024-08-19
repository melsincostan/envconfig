package envconfig

import (
	"fmt"
	"os"
	"strings"
)

const STRING_DEFAULT = ""

func getStringDef(env_name string, required bool, def_val string, has_default bool) (string, error) {
	def := STRING_DEFAULT
	if has_default {
		def = def_val
	}
	res, err := getString(env_name, true)
	if err != nil {
		if required {
			return STRING_DEFAULT, err
		}

		return def, nil
	}
	return res, nil
}

func getString(env_name string, required bool) (string, error) {
	raw := os.Getenv(env_name)
	val := strings.TrimSpace(raw)
	if required && len(val) < 1 {
		return "", fmt.Errorf("could not find %s even though it is required", env_name)
	}
	return val, nil
}
