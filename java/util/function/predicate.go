package function

/**
This is a port of java.util.function.Predicate.

* https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/function/Predicate.html
* https://github.com/openjdk/jdk/blob/jdk-21%2B35/src/java.base/share/classes/java/util/function/Predicate.java
*/

// Predicate is a type to represents a function that accepts one
// argument and produce one bool result an error.
type Predicate[T any] func(T) (bool, error)

type predicates[T any] []Predicate[T]

func newPredicates[T any](p1, p2 Predicate[T], p3 ...Predicate[T]) predicates[T] {
	if p1 == nil {
		return nil
	}
	if p2 == nil {
		return nil
	}
	if p3 != nil {
		for _, p := range p3 {
			if p == nil {
				return nil
			}
		}
	}
	ret := make(predicates[T], 0, 2+len(p3))
	ret = append(ret, p1, p2)
	ret = append(ret, p3...)
	return ret
}

// WrapNoErrPred adjusts a function that accepts one argument and
// produce one bool result to Predicate.
// If f is nil, thils function returns nil.
func WrapNoErrPred[T any](f func(T) bool) Predicate[T] {
	if f == nil {
		return nil
	}
	return Predicate[T](func(in T) (bool, error) { return f(in), nil })
}

// And returns a Predicate composed by arguments. The composed
// Predicate is a short-circuiting logical AND. The order of
// evaluating Predicates is as same as the order of arugments of this
// function.  If preceding Predicates return false or error, rest of
// Predicates are not evaluated.
func And[T any](p1, p2 Predicate[T], p3 ...Predicate[T]) Predicate[T] {
	ps := newPredicates(p1, p2, p3...)
	if ps == nil {
		return nil
	}
	return func(in T) (bool, error) {
		for _, p := range ps {
			ok, err := p(in)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		}
		return true, nil
	}
}

// Or returns a Predicate composed by arguments. The composed
// Predicate is a short-circuiting logical OR. Order of evaluating
// Predicates is as same as the order of arugments of this function.
// If preceding Predicates return true or error, rest of Predicates
// are not evaluated.
func Or[T any](p1, p2 Predicate[T], p3 ...Predicate[T]) Predicate[T] {
	ps := newPredicates(p1, p2, p3...)
	if ps == nil {
		return nil
	}
	return func(in T) (bool, error) {
		for _, p := range ps {
			ok, err := p(in)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	}
}

// Not returns a predicate which returns negation of the supplied
// predicate.
func Not[T any](t Predicate[T]) Predicate[T] {
	if t == nil {
		return nil
	}
	return func(in T) (bool, error) {
		ok, err := t(in)
		if err != nil {
			return false, err
		}
		return !ok, nil
	}
}

// ComparableEquals returns predicate which tests two comparable
// instance is same or not.
func ComparableEquals[T comparable](i T) Predicate[T] {
	return WrapNoErrPred(func(j T) bool { return i == j })
}
