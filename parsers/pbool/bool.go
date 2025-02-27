package pbool

import (
	"strings"

	"github.com/melsincostan/envconfig/parsers/pstring"
)

var TRUTHY_VALUES = map[string]bool{ // bool value doesn't get used
	"true": true,
	"yes":  true,
	"y":    true,
	"1":    true,
}

const BOOL_DEFAULT = "false"

func Parse(env_name string, required bool, def_val string, has_default bool) (bool, error) {
	def := truthy(BOOL_DEFAULT)
	if !required && has_default {
		def = truthy(def_val)
	}

	raw, err := pstring.Parse(env_name, true)

	if err != nil {
		if required {
			return truthy(BOOL_DEFAULT), err
		} else {
			return def, nil
		}
	}
	return truthy(strings.ToLower(raw)), nil
}

func truthy(s string) bool {
	_, ok := TRUTHY_VALUES[s]
	return ok
}
