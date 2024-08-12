package envconfig

import "strconv"

func getInt(env_name string, required bool) (int64, error) {
	raw, err := getString(env_name, true)
	if err != nil {
		if required {
			return 0, err
		} else {
			return 0, nil // can't parse an empty string, but it isn't required so return 0
		}
	}

	res, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}
