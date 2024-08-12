package envconfig

import (
	"fmt"
	"reflect"
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

		if !f_v.CanSet() {
			return nil, fmt.Errorf("Field %s: not assignable", f_t.Name)
		}

		switch f_t.Type.Kind() {
		default:
			return nil, fmt.Errorf("Field %s: unsupported type %s", f_t.Name, f_v.Kind())
		}
	}
	return ptr, nil
}
