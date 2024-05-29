package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/dairyo/j2g/java/util/function/consumer"
	"github.com/dairyo/j2g/java/util/function/function"
	"github.com/dairyo/j2g/java/util/function/predicate"
	"github.com/dairyo/j2g/java/util/function/supplier"
)

type optional[T any] interface {
	IfPresent(consumer.Consumer[T]) error
	Filter(predicate.Predicate[T]) *Optional[T]
	Or(supplier.Supplier[*Optional[T]]) *Optional[T]
	Get() T
	IsPresent() bool
	IsEmpty() bool
	Error() error
}

type Optional[T any] struct {
	o optional[T]
}

func (o *Optional[T]) IfPresent(c consumer.Consumer[T]) error {
	return o.o.IfPresent(c)
}

func (o *Optional[T]) Filter(p predicate.Predicate[T]) *Optional[T] {
	return o.o.Filter(p)
}

func (o *Optional[T]) Or(s supplier.Supplier[*Optional[T]]) *Optional[T] {
	return o.o.Or(s)
}

func (o *Optional[T]) Get() T {
	return o.o.Get()
}

func (o *Optional[T]) IsPresent() bool {
	return o.o.IsPresent()
}

func (o *Optional[T]) IsEmpty() bool {
	return o.o.IsEmpty()
}

func (o *Optional[T]) Error() error {
	return o.o.Error()
}

func isNilable(k reflect.Kind) bool {
	switch k {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice:
		return true
	default:
		return false
	}
}

func NewOptional[T any](v T) *Optional[T] {
	rv := reflect.ValueOf(v)
	vt := rv.Type()
	if isNilable(vt.Kind()) {
		if rv.IsNil() {
			return newErr[T](ErrEmpty)
		}
	}
	vo := &valueOptional[T]{val: v}
	ret := &Optional[T]{}
	vo.parent = ret
	ret.o = vo
	return ret
}

func Map[T, U any](v *Optional[T], f function.Function[T, U]) *Optional[U] {
	if v == nil {
		return newErr[U](ErrMapNilOptinal)
	}
	if f == nil {
		return newErr[U](ErrMapNilFunction)
	}

	switch val := v.o.(type) {
	case *errorOptional[T]:
		return newErr[U](fmt.Errorf("invalid optional is passed: %w", val.err))
	case *valueOptional[T]:
		ret, err := f(val.val)
		if err != nil {
			return newErr[U](fmt.Errorf("function returns err: %w", err))
		}
		return NewOptional(ret)
	default:
		return newErr[U](errors.New("unknown Optional type"))
	}
}

func FlatMap[T, U any](v *Optional[T], f function.Function[T, *Optional[U]]) *Optional[U] {
	if v == nil {
		return newErr[U](ErrMapNilOptinal)
	}
	if f == nil {
		return newErr[U](ErrMapNilFunction)
	}

	switch val := v.o.(type) {
	case *errorOptional[T]:
		return newErr[U](fmt.Errorf("invalid optional is passed: %w", val.err))
	case *valueOptional[T]:
		ret, err := f(val.val)
		if err != nil {
			return newErr[U](fmt.Errorf("function returns err: %w", err))
		}
		return ret
	default:
		return newErr[U](errors.New("unknown Optional type"))
	}
}
