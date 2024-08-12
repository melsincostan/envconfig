package envconfig

import "strconv"

func getInt(env_name string, required bool) (int, error) {
	raw, err := getString(env_name, true)
	if err != nil {
		if required {
			return 0, err
		} else {
			return 0, nil // can't parse an empty string, but it isn't required so return 0
		}
	}

	res, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	return res, nil
}
