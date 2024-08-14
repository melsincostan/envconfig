package envconfig

import (
	"fmt"
	"os"
	"testing"
)

func TestGetString(t *testing.T) {
	valid_env_name := "TEST_STRING_ENV_NAME"
	invalid_env_name := fmt.Sprintf("%s_INVALID", valid_env_name)

	valid_value := "TEST_STRING_ENV_VALUE"
	whitespace_value := "  \t\t\n\n\r\n\r\n"
	trimmable_value := fmt.Sprintf("%s%s", whitespace_value, valid_value)

	cases := []struct {
		Name        string
		EnvName     string
		EnvVal      string
		Required    bool
		Expect      string
		ExpectError bool
	}{
		{"EnvOK_ValOK_required", valid_env_name, valid_value, true, valid_value, false},
		{"EnvOK_ValOK_optional", valid_env_name, valid_value, false, valid_value, false},
		{"EnvBAD_ValOK_required", invalid_env_name, valid_value, true, "", true},
		{"EnvBAD_ValOK_optional", invalid_env_name, valid_value, false, "", false},
		{"EnvOK_ValBAD_required", valid_env_name, whitespace_value, true, "", true},
		{"EnvOK_ValBAD_optional", valid_env_name, whitespace_value, false, "", false},
		{"EnvBAD_ValBAD_required", invalid_env_name, whitespace_value, true, "", true},
		{"EnvBAD_ValBAD_optional", invalid_env_name, whitespace_value, false, "", false},
		{"EnvOK_ValTRIM_required", valid_env_name, trimmable_value, true, valid_value, false},
		{"EnvOK_ValTRIM_optional", valid_env_name, trimmable_value, false, valid_value, false},
		{"EnvBAD_ValTRIM_required", invalid_env_name, trimmable_value, true, "", true},
		{"EnvBAD_ValTRIM_optional", invalid_env_name, trimmable_value, false, "", false},
	}

	for _, c := range cases {
		c := c                             // local scope
		t.Run(c.Name, func(t *testing.T) { // can't run in parallel, since there might otherwise be race conditions with setting / unsetting env vars
			os.Setenv(c.EnvName, c.EnvVal)
			defer os.Unsetenv(c.EnvName) // ensure to clean up after ourselves :3

			res, err := getString(valid_env_name, c.Required)
			if err != nil && !c.ExpectError {
				t.Errorf("expected no error, got '%s'", err.Error())
			} else if err == nil && c.ExpectError {
				t.Error("expected to see an error")
			}

			if res != c.Expect {
				t.Errorf("wanted '%s', got '%s'", c.Expect, res)
			}

		})
	}
}
