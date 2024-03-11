package function

/**
This is a port of java.util.function.Function.

https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/function/Function.html
https://github.com/openjdk/jdk/blob/jdk-21%2B35/src/java.base/share/classes/java/util/function/Function.java
*/

import (
	"errors"
)

var (
	ErrNilFunc = errors.New("Sequence received nil Function")
)

// Function represents a function that accepts one argument and
// produces one result.
type Function[T any, U any] interface {
	// Apply applies this function to a given argument.
	Apply(T) (U, error)
}

// FunctionFunc is a type to adjust a function to Function interface.
type FunctionFunc[T any, U any] func(T) (U, error)

// Apply applies wrapped function to a given argument.
func (f FunctionFunc[T, U]) Apply(in T) (U, error) {
	return f(in)
}

// NewFunctionFunc generate Function from a function.
// This is utility function to infer types.
func NewFunctionFunc[T any, U any](f func(T) (U, error)) Function[T, U] {
	return FunctionFunc[T, U](f)
}

// NewFunctionFuncWithNoError generate function from a function which
// does not return error.
func NewFunctionFuncWithNoError[T any, U any](f func(T) U) Function[T, U] {
	return FunctionFunc[T, U](func(in T) (U, error) { return f(in), nil })
}

// Concatenate concatenates two functions.
// Returned value from first function becomes input to second function.
// Concatenate returns nil if one of or all of inputted functions are nil.
// This is a replacement of java Function's andThen and compose methods.
func Concatenate[T any, U any, V any](f1 Function[T, U], f2 Function[U, V]) Function[T, V] {
	if f1 == nil {
		return nil
	}
	if f2 == nil {
		return nil
	}
	return NewFunctionFunc(func(t T) (V, error) {
		u, err := f1.Apply(t)
		if err != nil {
			return *new(V), err
		}
		return f2.Apply(u)
	})
}

// Identity generate a function which always returns its input argument.
func Identity[T any]() Function[T, T] {
	return NewFunctionFunc(func(in T) (T, error) { return in, nil })
}
