package envconfig

import (
	"fmt"
	"testing"
)

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
