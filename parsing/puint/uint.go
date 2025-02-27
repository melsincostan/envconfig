package puint

import (
	"fmt"
	"strconv"

	"github.com/melsincostan/envconfig/parsing/pstring"
)

const UINT_DEFAULT = uint64(0)

func Parse(env_name string, required bool, def_val string, has_default bool) (uint64, error) {
	def := UINT_DEFAULT
	if !required && has_default {
		err := error(nil)
		def, err = strconv.ParseUint(def_val, 10, 64)
		if err != nil {
			return UINT_DEFAULT, fmt.Errorf("invalid default value: %s", err.Error())
		}
	}

	raw, err := pstring.Parse(env_name, true)
	if err != nil {
		if required {
			return UINT_DEFAULT, err
		} else {
			return def, nil // can't parse an empty string, but it isn't required so return 0
		}
	}

	res, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return UINT_DEFAULT, err
	}
	return res, nil
}
