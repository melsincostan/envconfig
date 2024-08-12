package envconfig

import (
	"fmt"
	"reflect"
)

func parse[T any]() (*T, error) {
	ptr := new(T)
	// ptr_t := reflect.TypeOf(ptr)
	ptr_v := reflect.ValueOf(ptr)
	// obj_t := ptr_t.Elem()
	obj_v := ptr_v.Elem()

	if obj_v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", obj_v.Kind())
	}
	return ptr, nil
}
