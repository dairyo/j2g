package util

import (
	"errors"

	"github.com/dairyo/j2g/java/lang/runnable"
	"github.com/dairyo/j2g/java/util/function/consumer"
	"github.com/dairyo/j2g/java/util/function/predicate"
	"github.com/dairyo/j2g/java/util/function/supplier"
)

type errorOptional[T any] struct {
	parent *Optional[T]
	err    error
}

func newErr[T any](err error) *Optional[T] {
	eo := &errorOptional[T]{err: err}
	ret := &Optional[T]{}
	eo.parent = ret
	ret.o = eo
	return ret
}

func (e *errorOptional[T]) IsPresent() bool {
	return false
}

func (e *errorOptional[T]) IsEmpty() bool {
	return true
}

func (e *errorOptional[T]) IfPresent(_ consumer.Consumer[T]) error {
	return ErrNoValue
}

func (e *errorOptional[T]) IfPresentOrElse(_ consumer.Consumer[T], r runnable.Runnable) error {
	if r == nil {
		return ErrInvalidUsed
	}
	err := r()
	if err != nil {
		return err
	}
	return nil
}

func (e *errorOptional[T]) Filter(_ predicate.Predicate[T]) *Optional[T] {
	return e.parent
}

func (e *errorOptional[T]) Error() error {
	return e.err
}

func (e *errorOptional[T]) Get() T {
	return *new(T)
}

func (e *errorOptional[T]) Or(s supplier.Supplier[*Optional[T]]) *Optional[T] {
	if s == nil {
		return newErr[T](errors.Join(ErrNilSupplier, e.err))
	}
	ret, err := s()
	if err != nil {
		return newErr[T](errors.Join(ErrSupplierErr, err, e.err))
	}
	return ret
}
