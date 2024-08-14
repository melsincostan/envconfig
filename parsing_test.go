package envconfig

import (
	"fmt"
	"math"
	"os"
	"testing"
)

type structUnexportedField struct {
	unexportedField string
}

type structUnsupportedType struct {
	UnsupportedField any
}

type structString struct {
	StringField         string `env:"TEST_STRING" binding:"required"`
	OptionalStringField string `env:"OPTIONAL_TEST_STRING"`
}

type structInt struct {
	Int64Field        int64 `env:"TEST_INT64" binding:"REQUIRED"`
	OptionalInt8Field int8  `env:"TEST_INT8"`
}

func TestParseNotStruct(t *testing.T) {
	res, err := parse[int]()
	if err == nil {
		t.Error("expected to see an error")
	}

	if res != nil {
		t.Error("expected a nil result")
	}
}

func TestParseUnexportedField(t *testing.T) {
	res, err := parse[structUnexportedField]()
	if err == nil {
		t.Error("expected to see an error")
	}

	if res != nil {
		t.Error("expected a nil result")
	}
}

func TestParseUnsupportedType(t *testing.T) {
	res, err := parse[structUnsupportedType]()
	if err == nil {
		t.Error("expected to see an error")
	}

	if res != nil {
		t.Error("expected a nil result")
	}
}

func TestParseString(t *testing.T) {
	test_string := "TEST_STRING_VALUE"
	optional_test_string := "OPTIONAL_TEST_STRING_VALUE"
	os.Setenv("TEST_STRING", test_string)
	os.Setenv("OPTIONAL_TEST_STRING", optional_test_string)

	res, err := parse[structString]()

	if err != nil {
		t.Errorf("expected to see no error, got %s", err.Error())
	}

	if res.StringField != test_string {
		t.Errorf("wanted '%s', got '%s'", test_string, res.StringField)
	}

	if res.OptionalStringField != optional_test_string {
		t.Errorf("wanted '%s', got '%s'", optional_test_string, res.OptionalStringField)
	}

	os.Unsetenv("OPTIONAL_TEST_STRING")

	res, err = parse[structString]()

	if err != nil {
		t.Errorf("expected no error, got %s", err.Error())
	}

	if res.StringField != test_string {
		t.Errorf("wanted '%s', got '%s'", test_string, res.StringField)
	}

	if res.OptionalStringField != "" {
		t.Errorf("wanted empty string, got '%s'", res.OptionalStringField)
	}

	os.Unsetenv("TEST_STRING")

	res, err = parse[structString]()

	if err == nil {
		t.Error("expected to see an error")
	}

	if res != nil {
		t.Error("expected a nil result")
	}
}

func TestParseInt(t *testing.T) {
	test_value := math.MaxInt32
	int8_test_value := int8(5)
	os.Setenv("TEST_INT64", fmt.Sprintf("%d", test_value))
	os.Setenv("TEST_INT8", fmt.Sprintf("%d", int8_test_value))

	res, err := parse[structInt]()

	if err != nil {
		t.Errorf("expected to see no error, got %s", err.Error())
	}

	if int(res.Int64Field) != test_value {
		t.Errorf("wanted '%d', got '%d'", test_value, res.Int64Field)
	}

	if res.OptionalInt8Field != int8_test_value {
		t.Errorf("wanted '%d', got '%d'", int8_test_value, res.OptionalInt8Field)
	}

	os.Setenv("TEST_INT8", fmt.Sprintf("%d", test_value))

	res, err = parse[structInt]()

	if err != nil {
		t.Errorf("expected to see no error, got %s", err.Error())
	}

	if int(res.Int64Field) != test_value {
		t.Errorf("wanted '%d', got '%d'", test_value, res.Int64Field)
	}

	if res.OptionalInt8Field != int8(test_value) {
		t.Errorf("wanted '%d', got '%d'", int8(test_value), res.OptionalInt8Field)
	}

	os.Unsetenv("TEST_INT8")

	res, err = parse[structInt]()

	if err != nil {
		t.Errorf("expected no error, got %s", err.Error())
	}

	if int(res.Int64Field) != test_value {
		t.Errorf("wanted '%d', got '%d'", test_value, res.Int64Field)
	}

	if res.OptionalInt8Field != 0 {
		t.Errorf("wanted '0', got '%d'", res.OptionalInt8Field)
	}

	os.Unsetenv("TEST_INT64")

	res, err = parse[structInt]()

	if err == nil {
		t.Error("expected to see an error")
	}

	if res != nil {
		t.Error("expected a nil result")
	}
}
