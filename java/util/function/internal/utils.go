package internal

import "reflect"

func Cast[T any, U any]() func(T) U {
	tt := reflect.TypeOf((*T)(nil)).Elem()
	ut := reflect.TypeOf((*U)(nil)).Elem()
	switch ut.Kind() {
	case reflect.Interface:
		if tt.Implements(ut) {
			return func(in T) U { return reflect.ValueOf(in).Interface().(U) }
		}
	default:
		if tt.ConvertibleTo(ut) {
			return func(in T) U { return reflect.ValueOf(in).Convert(ut).Interface().(U) }
		}
	}
	return nil
}
