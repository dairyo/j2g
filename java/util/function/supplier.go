package function

/**
This is a port of java.util.function.Supplier.

* https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/function/Supplier.html
* https://github.com/openjdk/jdk/blob/jdk-21%2B35/src/java.base/share/classes/java/util/function/Supplier.java
*/

// Supplier is a type to represents a function that accepts no
// argument nad produces one result and error.  The aim of this type
// is to supply a data. There is no requirement that a new or distinct
// result to be returned each time the Supplier is invoked.
type Supplier[T any] func() (T, error)

// WrapNoErrSupplier adjusts a function that accepts no argument and
// return a result to Supplier.
// If f is nil, this function returns nil.
func WrapNoErrSupplier[T any](f func() T) Supplier[T] {
	if f == nil {
		return nil
	}
	return func() (T, error) { return f(), nil }
}
