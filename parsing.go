package envconfig

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func parse[T any]() (*T, error) {
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

		env_name, has_env_tag := f_t.Tag.Lookup("env")

		if !has_env_tag {
			env_name = strings.ToUpper(f_t.Name)
		}

		def_val, has_default := f_t.Tag.Lookup("default")

		required := strings.ToLower(f_t.Tag.Get("binding")) == "required"

		if !f_v.CanSet() {
			return nil, fmt.Errorf("field %s: not assignable", f_t.Name)
		}

		switch f_v.Interface().(type) {
		case string:
			res, err := getString(env_name, required)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.SetString(res)
		case int, int8, int16, int32, int64:
			res, err := getInt(env_name, required, def_val, has_default)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.SetInt(res) // TODO; check if this truncates if assigning a number with higher bitsize to a field with smaller bitsize (for example in16-size number into int8)/
		case float32, float64:
			res, err := getFloat(env_name, required, def_val, has_default)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.SetFloat(res)
		case uint, uint8, uint16, uint32, uint64:
			res, err := getUint(env_name, required, def_val, has_default)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.SetUint(res)
		case time.Duration:
			res, err := getDuration(env_name, required)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.Set(reflect.ValueOf(res))
		default:
			return nil, fmt.Errorf("field %s: unsupported type %s", f_t.Name, f_v.Kind())
		}
	}
	return ptr, nil
}
