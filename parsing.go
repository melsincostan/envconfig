package envconfig

import (
	"fmt"
	"reflect"
	"strings"
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

		required := strings.ToLower(f_t.Tag.Get("binding")) == "required"

		if !f_v.CanSet() {
			return nil, fmt.Errorf("Field %s: not assignable", f_t.Name)
		}

		switch f_t.Type.Kind() {
		case reflect.String:
			res, err := getString(env_name, required)
			if err != nil {
				return nil, fmt.Errorf("field %s: %s", f_t.Name, err.Error())
			}
			f_v.SetString(res)
		default:
			return nil, fmt.Errorf("Field %s: unsupported type %s", f_t.Name, f_v.Kind())
		}
	}
	return ptr, nil
}
