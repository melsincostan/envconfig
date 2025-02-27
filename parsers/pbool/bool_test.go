package pbool

import (
	"fmt"
	"os"
	"testing"
)

var valid_env_name = "TEST_BOOL_VAL"
var invalid_env_name = fmt.Sprintf("%s_INCORRECT", valid_env_name)

type truthCase struct {
	Name   string
	Input  string
	Expect bool
}

func TestTruthy(t *testing.T) {
	cases := []truthCase{
		{"false_falsey", "false", false},
		{"no_falsey", "no", false},
		{"aabbcc_falsey", "aabbcc", false},
		{"empty_falsey", "", false},
	}
	for val := range TRUTHY_VALUES {
		cases = append(cases, truthCase{fmt.Sprintf("%s_truthy", val), val, true})
	}

	for _, c := range cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			res := truthy(c.Input)
			if res != c.Expect {
				t.Errorf("wanted %t, got %t", c.Expect, res)
			}
		})
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		Name        string
		EnvName     string
		EnvVal      string
		Required    bool
		Expect      bool
		ExpectError bool
		Default     string
		HasDefault  bool
	}{
		{"EnvOK_required", valid_env_name, "true", true, true, false, "", false},
		{"EnvOK_optional", valid_env_name, "false", false, false, false, "", false},
		{"EnvBAD_required", invalid_env_name, "true", true, false, true, "", false},
		{"EnvBAD_optional", invalid_env_name, "true", false, false, false, "", false},
		{"EnvOK_required_default", valid_env_name, "true", true, true, false, "false", true},
		{"EnvOK_optional_default", valid_env_name, "true", false, true, false, "false", true},
		{"EnvBAD_required_default", invalid_env_name, "true", true, false, true, "true", true},
		{"EnvBAD_optional_default", invalid_env_name, "false", false, true, false, "true", true},
	}

	for _, c := range cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			os.Setenv(c.EnvName, c.EnvVal)
			defer os.Unsetenv(c.EnvName)

			res, err := Parse(valid_env_name, c.Required, c.Default, c.HasDefault)
			if err != nil && !c.ExpectError {
				t.Errorf("expected no error, got %s", err.Error())
			} else if err == nil && c.ExpectError {
				t.Errorf("expected to see an error")
			}

			if res != c.Expect {
				t.Errorf("wanted '%t', got %t", res, c.Expect)
			}
		})
	}
}
