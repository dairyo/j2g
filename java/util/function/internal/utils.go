package internal

import (
	"reflect"

	"github.com/dairyo/j2g/java/util/function"
)

func cast[T any](i any) (T, error) {
	ret, ok := i.(T)
	if !ok {
		var zero T
		return zero, function.ErrFailToCast
	}
	return ret, nil
}

func Cast[T any, U any]() func(T) (U, error) {
	tt := reflect.TypeOf((*T)(nil)).Elem()
	ut := reflect.TypeOf((*U)(nil)).Elem()
	switch ut.Kind() {
	case reflect.Interface:
		if tt.Implements(ut) {
			return func(in T) (U, error) { return cast[U](reflect.ValueOf(in).Interface()) }
		}
	default:
		if tt.ConvertibleTo(ut) {
			return func(in T) (U, error) { return cast[U](reflect.ValueOf(in).Convert(ut).Interface()) }
		}
		if tt.Kind() == reflect.Interface && ut.Implements(tt) {
			return func(in T) (U, error) { return cast[U](reflect.ValueOf(in).Interface()) }
		}
	}
	return nil
}
