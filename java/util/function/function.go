package function

/**
This is a port of java.util.function.Function.

* https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/function/Function.html
* https://github.com/openjdk/jdk/blob/jdk-21%2B35/src/java.base/share/classes/java/util/function/Function.java
*/

// Function is a type to represents a function that accepts one
// argument and produce one result and error.
type Function[T any, U any] func(T) (U, error)

// WrapNoErrFunc adjusts a function that accepts one argument and
// produce one result to Function.
func WrapNoErrFunc[T any, U any](f func(in T) U) Function[T, U] {
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
