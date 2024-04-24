package util

import (
	"errors"
	"fmt"
	"reflect"
)

/**
This is a port of java.util.Optional.

* https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/Optional.html
*https://github.com/openjdk/jdk/blob/jdk-21%2B35/src/java.base/share/classes/java/util/Optional.java#L86
*/

var (
	ErrArgumentIsNil = errors.New("an argument is nil")
)

type Optional[T any] struct {
	val T
}

func Empty[T any]() *Optional[T] {
	return &Optional[T]{}
}

func newOptional[T any](val T) *Optional[T] {
	return &Optional[T]{val: val}
}

func isNilable[T any](in T) bool {
	rt := reflect.TypeOf(in)
	switch rt.Kind() {
	case reflect.Array:
		return true
	case reflect.Chan:
		return true
	case reflect.Func:
		return true
	case reflect.Interface:
		return true
	case reflect.Pointer:
		return true
	case reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

func Of[T comparable](val T) (*Optional[T], error) {
	if isNilable(val) {
		if val == *new(T) {
			return nil, fmt.Errorf("value is nil for type %T\n%w", val, ErrArgumentIsNil)
		}
	}
	return newOptional(val), nil
}

func OfNilable[T comparable](val T) *Optional[T] {
	return newOptional(val)
}
