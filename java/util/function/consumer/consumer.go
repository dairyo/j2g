package consumer

/**
This is a port of java.util.function.Consumer.

* https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/function/Consumer.html
* https://github.com/openjdk/jdk/blob/jdk-21%2B35/src/java.base/share/classes/java/util/function/Consumer.java
*/

// Consumer is a type to represents a function that accepts one
// argument and returns error. Unlike other functional types, Consumer
// is expected to operate via side-effect.
type Consumer[T any] func(in T) error

// WrapNoErr adjusts a vvunction that accepts one argument and
// no return to Consumer.
// If f is nil, this function returns nil.
func WrapNoErr[T any](f func(in T)) Consumer[T] {
	if f == nil {
		return nil
	}
	return func(in T) error {
		f(in)
		return nil
	}
}

// Compose returns a Consumer composed by arguments. The composed
// Consumer evaluates Consumers passed as arguments. The order of
// evaluating Consumers is as the same as the order of arguments.  If
// preceding Consumers return error, rest of the Consumers are not
// evaluated.
func Compose[T any](c1 Consumer[T], c2 ...Consumer[T]) Consumer[T] {
	if c1 == nil {
		return nil
	}
	for _, c := range c2 {
		if c == nil {
			return nil
		}
	}
	return func(in T) error {
		if err := c1(in); err != nil {
			return err
		}
		for _, c := range c2 {
			if err := c(in); err != nil {
				return err
			}
		}
		return nil
	}
}
