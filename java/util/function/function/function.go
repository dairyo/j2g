package function

import (
	"fmt"

	"github.com/dairyo/j2g/java/util/function/internal"
)

/**
This is a port of java.util.function.Function.

* https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/function/Function.html
* https://github.com/openjdk/jdk/blob/jdk-21%2B35/src/java.base/share/classes/java/util/function/Function.java
*/

// Function is a type to represents a function that accepts one
// argument and produce one result and error.
type Function[T any, U any] func(T) (U, error)

// WrapNoErr adjusts a function that accepts one argument and produce
// one result to Function.
// If f is nil, this function returns nil.
func WrapNoErr[T any, U any](f func(in T) U) Function[T, U] {
	if f == nil {
		return nil
	}
	return func(in T) (U, error) { return f(in), nil }
}

// Compose composes two functions.
// Returned value from first function becomes input to second function.
// Compose returns nil if one of or both of inputted functions are nil.
// This is a replacement of java Function's andThen and compose methods.
func Compose[T any, U any, V any](f1 Function[T, U], f2 Function[U, V]) Function[T, V] {
	if f1 == nil {
		return nil
	}
	if f2 == nil {
		return nil
	}
	return (func(t T) (V, error) {
		u, err := f1(t)
		if err != nil {
			return *new(V), err
		}
		return f2(u)
	})
}

// Identity generate a function which always returns its input argument.
func Identity[T any]() Function[T, T] {
	return func(in T) (T, error) { return in, nil }
}

// Adjust adjusts a function to other function.
//
// Adjust is mainly used in arguments of [And] and [Or]. For example:
//
//  	f1 := func(b *bytes.Buffer) (*bytes.Buffer, error) {
//  		return b, nil
//  	}
//  	f2 := func(w io.Writer) (io.Writer, error) {
//  		return w, nil
//  	}
//  	b := &bytes.Buffer{}
//  	Compose(f1, Adjust[*bytes.Buffer, io.Writer, *bytes.Buffer, io.Writer](f2))(b)
//
// If U1 is an interface, T1 must implements U1. If U1 is a type, T1
// must be convertible to U1.  If U2 is an interface, T2 must
// implements U2. If U2 is a type, T2 must be convertible to U2.
//
// This function might panic. We recommend you should write adjusting
// function by your own like following:
//
//  	f1 := func(b *bytes.Buffer) (*bytes.Buffer, error) {
//  		return b, nil
//  	}
//  	f2 := func(w io.Writer) (io.Writer, error) {
//  		return w, nil
//  	}
//  	b := &bytes.Buffer{}
//  	Compose(
//  		f1,
//  		func(b *bytes.Buffer) (*bytes.Buffer, error) {
//  			r1, err := f2(b)
//  			if err != nil {
//  				return nil, err
//  			}
//  			r2, ok := r1.(*bytes.Buffer)
//  			if !ok {
//  				return nil, errors.New("fail to cast")
//  			}
//  			return r2, nil
//  		},
//  	)(b)
func Adjust[T1, U1, T2, U2 any](f func(U1) (U2, error)) func(T1) (T2, error) {
	if f == nil {
		return nil
	}
	cf1 := internal.Cast[T1, U1]()
	if cf1 == nil {
		return nil
	}
	cf2 := internal.Cast[U2, T2]()
	if cf2 == nil {
		return nil
	}
	return func(in T1) (T2, error) {
		u1, err := cf1(in)
		if err != nil {
			return *new(T2), fmt.Errorf("fail to cast argument from %T to %T: %w", in, u1, err)
		}
		ret, err := f(u1)
		if err != nil {
			return *new(T2), err
		}
		u2, err := cf2(ret)
		if err != nil {
			return *new(T2), fmt.Errorf("fail to cast return value from %T to %T: %w", ret, u2, err)
		}
		return u2, nil
	}
}
