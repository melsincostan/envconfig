package envconfig

import (
	"fmt"
	"os"
	"testing"
)

func TestGetUint(t *testing.T) {
	valid_env_name := "TEST_UINT_ENV_NAME"
	invalid_env_name := fmt.Sprintf("%s_INVALID", valid_env_name)

	valid_value := "5"
	valid_value_uint64 := uint64(5)
	whitespace_value := "  \t\t\n\n\r\n\r\n"
	trimmable_value := fmt.Sprintf("%s%s", whitespace_value, valid_value)
	invalid_value := "aabb"
	valid_negative_value := "-6"

	cases := []struct {
		Name        string
		EnvName     string
		EnvVal      string
		Required    bool
		Expect      uint64
		ExpectError bool
	}{
		{"EnvOK_ValOK_required", valid_env_name, valid_value, true, valid_value_uint64, false},
		{"EnvOK_ValOK_optional", valid_env_name, valid_value, false, valid_value_uint64, false},
		{"EnvBAD_ValOK_required", invalid_env_name, valid_value, true, 0, true},
		{"EnvBAD_ValOK_optional", invalid_env_name, valid_value, false, 0, false},
		{"EnvOK_ValWHITESPACE_required", valid_env_name, whitespace_value, true, 0, true},
		{"EnvOK_ValWHITESPACE_optional", valid_env_name, whitespace_value, false, 0, false},
		{"EnvBAD_ValWHITESPACE_required", invalid_env_name, whitespace_value, true, 0, true},
		{"EnvBAD_ValWHITESPACE_optional", invalid_env_name, whitespace_value, false, 0, false},
		{"EnvOK_ValTRIM_required", valid_env_name, trimmable_value, true, valid_value_uint64, false},
		{"EnvOK_ValTRIM_optional", valid_env_name, trimmable_value, false, valid_value_uint64, false},
		{"EnvBAD_ValTRIM_required", invalid_env_name, trimmable_value, true, 0, true},
		{"EnvBAD_ValTRIM_optional", invalid_env_name, trimmable_value, false, 0, false},
		{"EnvOK_ValBAD_required", valid_env_name, invalid_value, true, 0, true},
		{"EnvOK_ValBAD_optional", valid_env_name, invalid_value, false, 0, true},
		{"EnvBAD_ValBAD_required", invalid_env_name, invalid_value, true, 0, true},
		{"EnvBAD_ValBAD_optional", invalid_env_name, invalid_value, false, 0, false},
		{"EnvOK_ValNEG_required", valid_env_name, valid_negative_value, true, 0, true},
		{"EnvOK_ValNEG_optional", valid_env_name, valid_negative_value, false, 0, true},
		{"EnvBAD_ValNEG_required", invalid_env_name, valid_negative_value, true, 0, true},
		{"EnvBAD_ValNEG_optional", invalid_env_name, valid_negative_value, false, 0, false},
	}

	for _, c := range cases {
		c := c                             // local scope
		t.Run(c.Name, func(t *testing.T) { // can't run in parallel, since there might otherwise be race conditions with setting / unsetting env vars
			os.Setenv(c.EnvName, c.EnvVal)
			defer os.Unsetenv(c.EnvName) // ensure to clean up after ourselves :3

			res, err := getUint(valid_env_name, c.Required)
			if err != nil && !c.ExpectError {
				t.Errorf("expected no error, got '%s'", err.Error())
			} else if err == nil && c.ExpectError {
				t.Error("expected to see an error")
			}

			if res != c.Expect {
				t.Errorf("wanted '%d', got '%d'", c.Expect, res)
			}

		})
	}
}
