package util

import (
	"errors"

	"github.com/dairyo/j2g/java/lang/runnable"
	"github.com/dairyo/j2g/java/util/function/consumer"
	"github.com/dairyo/j2g/java/util/function/predicate"
	"github.com/dairyo/j2g/java/util/function/supplier"
)

type valueOptional[T any] struct {
	parent *Optional[T]
	val    T
}

func (o *valueOptional[T]) IsPresent() bool {
	return true
}

func (o *valueOptional[T]) IsEmpty() bool {
	return false
}

func (o *valueOptional[T]) IfPresent(c consumer.Consumer[T]) error {
	if c == nil {
		return ErrNilConsumer
	}
	return c(o.val)
}

func (o *valueOptional[T]) IfPresentOrElse(c consumer.Consumer[T], _ runnable.Runnable) error {
	return o.IfPresent(c)
}

func (o *valueOptional[T]) Filter(p predicate.Predicate[T]) *Optional[T] {
	if p == nil {
		return newErr[T](ErrNilPredicate)
	}
	ok, err := p(o.val)
	if !ok {
		return newErr[T](ErrPredicateFailed)
	}
	if err != nil {
		return newErr[T](errors.Join(ErrPredicateErr, err))
	}
	return o.parent
}

func (o *valueOptional[T]) Or(s supplier.Supplier[*Optional[T]]) *Optional[T] {
	if s == nil {
		return newErr[T](ErrNilSupplier)
	}
	return o.parent
}

func (o *valueOptional[T]) Get() (T, error) {
	return o.val, nil
}

func (o *valueOptional[T]) Error() error {
	return nil
}
