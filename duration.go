package envconfig

import (
	"fmt"
	"time"
)

const DURATION_DEFAULT = 0 * time.Second

func getDuration(env_name string, required bool, def_val string, has_default bool) (time.Duration, error) {
	def := DURATION_DEFAULT
	if !required && has_default {
		err := error(nil)
		def, err = time.ParseDuration(def_val)
		if err != nil {
			return DURATION_DEFAULT, fmt.Errorf("invalid default value: %s", err.Error())
		}
	}

	raw, err := getString(env_name, true)
	if err != nil {
		if required {
			return DURATION_DEFAULT, err
		} else {
			return def, nil // can't parse an empty string, but it isn't required so return 0
		}
	}

	res, err := time.ParseDuration(raw)
	if err != nil {
		return DURATION_DEFAULT, err
	}
	return res, nil
}
