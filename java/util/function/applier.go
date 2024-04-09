package function

/**
This is a port of java.util.function.Function.

* https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/function/Function.html
* https://github.com/openjdk/jdk/blob/jdk-21%2B35/src/java.base/share/classes/java/util/function/Function.java
*/

import (
	"errors"
)

var (
	ErrNilFunc = errors.New("Sequence received nil Function")
)

// Applier applies a function that accepts one argument and
// produces one result.
type Applier[T any, U any] interface {
	// Apply applies this function to a given argument.
	Apply(T) (U, error)
}

// ApplierFunc is a type to adjust a function to Applier interface.
type ApplierFunc[T any, U any] func(T) (U, error)

// Apply applies wrapped function to a given argument.
func (f ApplierFunc[T, U]) Apply(in T) (U, error) {
	return f(in)
}

// NewApplierFunc generate Applier from a function.
// This is utility function to infer types.
func NewApplierFunc[T any, U any](f func(T) (U, error)) Applier[T, U] {
	return ApplierFunc[T, U](f)
}

// NewApplierFuncWithNoError generate function from a function which
// does not return error.
func NewApplierFuncWithNoError[T any, U any](f func(T) U) Applier[T, U] {
	return ApplierFunc[T, U](func(in T) (U, error) { return f(in), nil })
}

// Compose composes two functions.
// Returned value from first function becomes input to second function.
// Compose returns nil if one of or both of inputted functions are nil.
// This is a replacement of java Function's andThen and compose methods.
func Compose[T any, U any, V any](f1 Applier[T, U], f2 Applier[U, V]) Applier[T, V] {
	if f1 == nil {
		return nil
	}
	if f2 == nil {
		return nil
	}
	return NewApplierFunc(func(t T) (V, error) {
		u, err := f1.Apply(t)
		if err != nil {
			return *new(V), err
		}
		return f2.Apply(u)
	})
}

// Identity generate a function which always returns its input argument.
func Identity[T any]() Applier[T, T] {
	return NewApplierFunc(func(in T) (T, error) { return in, nil })
}
