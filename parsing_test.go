package envconfig

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"
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
	UnusedInt16Field  int16
	UnusedInt32Field  int32
}

type structFloat struct {
	Float64Field         float64 `env:"TEST_FLOAT64" binding:"required"`
	OptionalFloat32Field float32 `env:"TEST_FLOAT32"`
}

type structUint struct {
	Uint64Field        uint64 `env:"TEST_UINT64" binding:"REQUIRED"`
	OptionalUint8Field uint8  `env:"TEST_UINT8"`
	UnusedUint16Field  uint16
	UnusedUint32Field  uint32
}

type structDuration struct {
	DurationField         time.Duration `env:"TEST_DURATION" binding:"required"`
	OptionalDurationField time.Duration `env:"OPTIONAL_TEST_DURATION"`
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
	test_value := int64(math.MaxInt32)
	int8_test_value := int8(5)
	os.Setenv("TEST_INT64", fmt.Sprintf("%d", test_value))
	os.Setenv("TEST_INT8", fmt.Sprintf("%d", int8_test_value))

	res, err := parse[structInt]()

	if err != nil {
		t.Errorf("expected to see no error, got %s", err.Error())
	}

	if res.Int64Field != test_value {
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

	if res.Int64Field != test_value {
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

	if res.Int64Field != test_value {
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

func TestParseFloat(t *testing.T) {
	test_value := float64(7.3)
	float32_test_value := float32(6.4)
	delta := 0.00001
	os.Setenv("TEST_FLOAT64", fmt.Sprintf("%.2f", test_value))
	os.Setenv("TEST_FLOAT32", fmt.Sprintf("%.2f", float32_test_value))

	res, err := parse[structFloat]()

	if err != nil {
		t.Errorf("expected to see no error, got %s", err.Error())
	}

	if res.Float64Field > test_value+delta || res.Float64Field < test_value-delta {
		t.Errorf("wanted '%.2f', got '%.2f'", test_value, res.Float64Field)
	}

	if res.OptionalFloat32Field > float32_test_value+float32(delta) || res.OptionalFloat32Field < float32_test_value-float32(delta) {
		t.Errorf("wanted '%.2f', got '%.2f'", float32_test_value, res.OptionalFloat32Field)
	}

	os.Unsetenv("TEST_FLOAT32")

	res, err = parse[structFloat]()

	if err != nil {
		t.Errorf("expected no error, got %s", err.Error())
	}

	if res.Float64Field > test_value+delta || res.Float64Field < test_value-delta {
		t.Errorf("wanted '%.2f', got '%.2f'", test_value, res.Float64Field)
	}

	if res.OptionalFloat32Field != 0 {
		t.Errorf("wanted '0', got '%.2f'", res.OptionalFloat32Field)
	}

	os.Unsetenv("TEST_FLOAT64")

	res, err = parse[structFloat]()

	if err == nil {
		t.Error("expected to see an error")
	}

	if res != nil {
		t.Error("expected a nil result")
	}
}

func TestParseUint(t *testing.T) {
	test_value := uint64(math.MaxInt32)
	uint8_test_value := uint8(5)
	os.Setenv("TEST_UINT64", fmt.Sprintf("%d", test_value))
	os.Setenv("TEST_UINT8", fmt.Sprintf("%d", uint8_test_value))

	res, err := parse[structUint]()

	if err != nil {
		t.Errorf("expected to see no error, got %s", err.Error())
	}

	if res.Uint64Field != test_value {
		t.Errorf("wanted '%d', got '%d'", test_value, res.Uint64Field)
	}

	if res.OptionalUint8Field != uint8_test_value {
		t.Errorf("wanted '%d', got '%d'", uint8_test_value, res.OptionalUint8Field)
	}

	os.Setenv("TEST_UINT8", fmt.Sprintf("%d", test_value))

	res, err = parse[structUint]()

	if err != nil {
		t.Errorf("expected to see no error, got %s", err.Error())
	}

	if res.Uint64Field != test_value {
		t.Errorf("wanted '%d', got '%d'", test_value, res.Uint64Field)
	}

	if res.OptionalUint8Field != uint8(test_value) {
		t.Errorf("wanted '%d', got '%d'", uint8(test_value), res.OptionalUint8Field)
	}

	os.Unsetenv("TEST_UINT8")

	res, err = parse[structUint]()

	if err != nil {
		t.Errorf("expected no error, got %s", err.Error())
	}

	if res.Uint64Field != test_value {
		t.Errorf("wanted '%d', got '%d'", test_value, res.Uint64Field)
	}

	if res.OptionalUint8Field != 0 {
		t.Errorf("wanted '0', got '%d'", res.OptionalUint8Field)
	}

	os.Unsetenv("TEST_UINT64")

	res, err = parse[structUint]()

	if err == nil {
		t.Error("expected to see an error")
	}

	if res != nil {
		t.Error("expected a nil result")
	}
}

func TestParseDuration(t *testing.T) {
	test_duration := 5 * time.Minute
	optional_test_duration := 6*time.Hour + 30*time.Minute
	os.Setenv("TEST_DURATION", fmt.Sprintf("%s", test_duration.String()))
	os.Setenv("OPTIONAL_TEST_DURATION", fmt.Sprintf("%s", optional_test_duration.String()))

	res, err := parse[structDuration]()

	if err != nil {
		t.Errorf("expected to see no error, got %s", err.Error())
	}

	if res.DurationField != test_duration {
		t.Errorf("wanted '%s', got '%s'", test_duration, res.DurationField)
	}

	if res.OptionalDurationField != optional_test_duration {
		t.Errorf("wanted '%s', got '%s'", optional_test_duration, res.OptionalDurationField)
	}

	os.Unsetenv("OPTIONAL_TEST_DURATION")

	res, err = parse[structDuration]()

	if err != nil {
		t.Errorf("expected no error, got %s", err.Error())
	}

	if res.DurationField != test_duration {
		t.Errorf("wanted '%s', got '%s'", test_duration, res.DurationField)
	}

	if res.OptionalDurationField != 0 {
		t.Errorf("wanted empty string, got '%s'", res.OptionalDurationField)
	}

	os.Unsetenv("TEST_DURATION")

	res, err = parse[structDuration]()

	if err == nil {
		t.Error("expected to see an error")
	}

	if res != nil {
		t.Error("expected a nil result")
	}
}
