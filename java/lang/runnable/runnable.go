package runnable

/**
This is a port of java.lang.Runnable.

* https://docs.oracle.com/javase/jp/8/docs/api/java/lang/Runnable.html
* https://github.com/openjdk/jdk/blob/jdk-21%2B35/src/java.base/share/classes/java/lang/Runnable.java
*/

// Runable represents an function that does not return result but may
// return error.
type Runnable func() error

// WrapNoErr adjusts a function that does not return result to
// [Runnable].
func WrapNoErr(f func()) Runnable {
	if f == nil {
		return nil
	}
	return func() error {
		f()
		return nil
	}
}
