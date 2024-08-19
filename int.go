package envconfig

import (
	"fmt"
	"strconv"
)

const INT_DEFAULT = int64(0)

func getInt(env_name string, required bool, def_val string, has_default bool) (int64, error) {
	def := INT_DEFAULT
	if !required && has_default {
		err := error(nil)
		def, err = strconv.ParseInt(def_val, 10, 64)
		if err != nil {
			return INT_DEFAULT, fmt.Errorf("unparseable default value: %s", err.Error())
		}
	}

	raw, err := getString(env_name, true)
	if err != nil {
		if required {
			return INT_DEFAULT, err
		} else {
			return def, nil // can't parse an empty string, but it isn't required so return 0
		}
	}

	res, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return INT_DEFAULT, err
	}
	return res, nil
}
