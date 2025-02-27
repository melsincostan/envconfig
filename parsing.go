package envconfig

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/melsincostan/envconfig/parsing/pbool"
	"github.com/melsincostan/envconfig/parsing/pduration"
	"github.com/melsincostan/envconfig/parsing/pfloat"
	"github.com/melsincostan/envconfig/parsing/pint"
	"github.com/melsincostan/envconfig/parsing/pstring"
	"github.com/melsincostan/envconfig/parsing/puint"
)

// Parse creates a struct of type T and attempts to fill it using environment variables.
// This behaviour can be adjusted through struct tags.
// The "env" struct tag will set the name of the environment variable this function looks for.
// If strings are passed as arguments to Parse, then they will be joined with "_" and the resulting string, with a trailing "_", will be used as a prefix for the capitalized field name or the value of the "env" tag.
// If it isn't set, then the uppercased name of the field will be used.
// If the "binding" tag is set to "required", then an error will be thrown if the environment variable is unset.
// Otherwise, a default value will be used.
// The default value can be set by using the "default" tag.
func Parse[T any](scopes ...string) (*T, error) {
	ptr := new(T)
	ptr_t := reflect.TypeOf(ptr)
	ptr_v := reflect.ValueOf(ptr)
	obj_t := ptr_t.Elem()
	obj_v := ptr_v.Elem()

	if obj_v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", obj_v.Kind())
	}

	for i := 0; i < obj_v.NumField(); i++ {
		f_t := obj_t.Field(i)
		f_v := obj_v.Field(i)

		raw_env_name, has_env_tag := f_t.Tag.Lookup("env")
		env_name := name(raw_env_name, scopes...)

		if raw_env_name == "-" {
			continue
		}

		if !has_env_tag {
			env_name = name(strings.ToUpper(f_t.Name), scopes...)
		}

		def_val, has_default := f_t.Tag.Lookup("default")

		required := strings.ToLower(f_t.Tag.Get("binding")) == "required"

		if !f_v.CanSet() {
			return nil, fmt.Errorf("field %s: not assignable", f_t.Name)
		}

		switch f_v.Interface().(type) {
		case string:
			res, err := pstring.ParseWithDefault(env_name, required, def_val, has_default)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.SetString(res)
		case int, int8, int16, int32, int64:
			res, err := pint.Parse(env_name, required, def_val, has_default)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.SetInt(res) // TODO; check if this truncates if assigning a number with higher bitsize to a field with smaller bitsize (for example in16-size number into int8)/
		case float32, float64:
			res, err := pfloat.Parse(env_name, required, def_val, has_default)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.SetFloat(res)
		case uint, uint8, uint16, uint32, uint64:
			res, err := puint.Parse(env_name, required, def_val, has_default)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.SetUint(res)
		case time.Duration:
			res, err := pduration.Parse(env_name, required, def_val, has_default)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.Set(reflect.ValueOf(res))
		case bool:
			res, err := pbool.Parse(env_name, required, def_val, has_default)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.SetBool(res)
		default:
			return nil, fmt.Errorf("field %s: unsupported type %s", f_t.Name, f_v.Kind())
		}
	}
	return ptr, nil
}

func name(name string, parts ...string) string {
	if len(parts) > 0 {
		return fmt.Sprintf("%s_%s", strings.Join(parts, "_"), name)
	}
	return name
}
