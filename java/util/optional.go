// This is a port of java.util.Optional.
//
//   - https://docs.oracle.com/javase/jp/21/docs/api/java.base/java/util/Optional.html
//   - https://github.com/openjdk/jdk/blob/jdk-21%2B35/src/java.base/share/classes/java/util/Optional.java
package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/dairyo/j2g/java/lang/runnable"
	"github.com/dairyo/j2g/java/util/function/consumer"
	"github.com/dairyo/j2g/java/util/function/function"
	"github.com/dairyo/j2g/java/util/function/predicate"
	"github.com/dairyo/j2g/java/util/function/supplier"
)

type optional[T any] interface {
	IfPresent(consumer.Consumer[T]) error
	IfPresentOrElse(consumer.Consumer[T], runnable.Runnable) error
	Filter(predicate.Predicate[T]) *Optional[T]
	Or(supplier.Supplier[*Optional[T]]) *Optional[T]
	Get() (T, error)
	IsPresent() bool
	IsEmpty() bool
	Error() error
}

// Optional is a container of a value.
type Optional[T any] struct {
	o optional[T]
}

// IfPresent executes [consumer.Consumer] c if [Optional.IsPresent]
// returns true. IfPresent may return following errors:
//   - [ErrNoValue] is returned if [Optional.IsPresent] returns false.
//   - [ErrNilConsumer] is returned if c is nil if [Optional.IsPresent]
//   returns true and c is nil.
//   - error returned by c is returned if c returns error.
func (o *Optional[T]) IfPresent(c consumer.Consumer[T]) error {
	return o.o.IfPresent(c)
}

// IfPresentOrElse executes [consumer.Consumer] c if
// [Optional.IsPresent] returns true. It executes runnable.Runnable r
// if [Optional.IsPresent] returns false. IfPresent may return
// following errors:
//   - [ErrNilRunnable] is returned if [Optional.IsPresent] is false and
//   r is nil.
//   - [ErrNilConsumer] is returned if c is nil if [Optional.IsPresent]
//   is true and c is nil.
//   - error returned by r is returned if [Optional.IsPresent] is
//   false and r returns error.
//   - error returned by c is returned if [Optional.IsPresent] is
//   true and c returns error.
func (o *Optional[T]) IfPresentOrElse(c consumer.Consumer[T], r runnable.Runnable) error {
	return o.o.IfPresentOrElse(c, r)
}

// Filter returns this Optional instance if the value of this Optional
// matches [predicate.Predicate] p. If the value of thie Optional does
// not match p or the value is empty, return an Optional instnce which
// does not have value. If empty Optional is returned, you can get the
// reason with [Optional.Error]. In this case [Optional.Error] may
// return following errors:
//   - [ErrNilPredicate] is returned if [predicate.Predicate] p is nil.
//   - [ErrPredicateFailed] is returned if [predicate.Predicate] p returns false.
//   - [ErrPredicateErr] is returned if [predicate.Predicate] p
//   returns error. In this case, the error returned by P is joined
//   with [errors.Join].
//   - If Optional instance calling Filter is already empty,
//   [Optional.Error] returns the original error.
func (o *Optional[T]) Filter(p predicate.Predicate[T]) *Optional[T] {
	return o.o.Filter(p)
}

// Or returns this Optional instance if the value of this Optional is
// present. Otherwise returns an Optional produced by
// [supplier.Supplier] s. Or returns empty Optional in following
// cases:
//   - [supplier.Supplier] s is Nil. In this case, [Optional.Error] returns [ErrNilSupplier]
//   - [supplier.Supplier] s returns error. In this case [Optional.Error] returns an error which contains [ErrSupplierErr], error returned by [supplier.Supplier] and error which is contained by base empty Optional.
func (o *Optional[T]) Or(s supplier.Supplier[*Optional[T]]) *Optional[T] {
	return o.o.Or(s)
}

// Get returns a value in this Optional instance if
// [Optional.IsPresent] is true. Otherwise Get returns [ErrNoValue].
func (o *Optional[T]) Get() (T, error) {
	return o.o.Get()
}

// IsPresent returns true if this Optional instance has a
// value. Otherwise return false.
func (o *Optional[T]) IsPresent() bool {
	return o.o.IsPresent()
}

// IsEmpty returns true if this Optional instance does not have a
// value. Otherwise return true.
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

// NewOptional returns an [Optional] instance holding value v.
// If v is nil, NewOptional returns empty [Optional].
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

// Map returns new [Optional] instance holding the result of applying
// the given mapping function f.
// If v is empty or nil NewOptional returns empty [Optional]
// instance. NewOptional also returns empty [Optional] instance if f
// is nil.
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
