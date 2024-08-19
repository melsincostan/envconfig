package envconfig

import (
	"fmt"
	"strconv"
)

const FLOAT_DEFAULT = float64(0)

func getFloat(env_name string, required bool, def_val string, has_default bool) (float64, error) {
	def := FLOAT_DEFAULT
	if !required && has_default {
		err := error(nil)
		def, err = strconv.ParseFloat(def_val, 64)
		if err != nil {
			return FLOAT_DEFAULT, fmt.Errorf("invalid default value: %s", err.Error())
		}
	}
	raw, err := getString(env_name, true)
	if err != nil {
		if required {
			return FLOAT_DEFAULT, err
		} else {
			return def, nil // can't parse an empty string, but it isn't required so return 0
		}
	}

	res, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return FLOAT_DEFAULT, err
	}
	return res, nil
}
