package predicate

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func checkPredicate[T any](t testing.TB, f Predicate[T], in T, want bool) {
	t.Helper()
	got, err := f(in)
	if err != nil {
		t.Errorf("must not return error but %q", err.Error())
		return
	}
	if got != want {
		t.Errorf("must return %t but %t", want, got)
		return
	}
}

func checkPredicateError[T any](t testing.TB, f Predicate[T], in T, want error) {
	t.Helper()
	v, got := f(in)
	if v == true {
		t.Error("must return false but true")
	}
	if diff := cmp.Diff(got.Error(), want.Error()); diff != "" {
		t.Error(diff)
		return
	}
}

func wnep[T any](f func(T) bool) Predicate[T] {
	return WrapNoErr(f)
}

func and[T any](t1, t2 Predicate[T], t3 ...Predicate[T]) Predicate[T] {
	return And(t1, t2, t3...)
}

func or[T any](t1, t2 Predicate[T], t3 ...Predicate[T]) Predicate[T] {
	return Or(t1, t2, t3...)
}

func TestWrapNoErr(t *testing.T) {
	f1 := wnep(func(in int) bool { return in == 1 })
	checkPredicate(t, f1, 1, true)
	checkPredicate(t, f1, 2, false)

	f2 := wnep((func(int) bool)(nil))
	if f2 != nil {
		t.Error("f2 must be nil")
	}
}

func TestAnd(t *testing.T) {
	f := and(
		wnep(func(in int) bool { return in%2 == 0 }),
		wnep(func(in int) bool { return in%5 == 0 }),
	)
	checkPredicate(t, f, 10, true)
	checkPredicate(t, f, 2, false)
	checkPredicate(t, f, 5, false)
	checkPredicate(t, f, 3, false)

	f2 := and(
		wnep(func(in int) bool { return in%2 == 0 }),
		wnep(func(in int) bool { return in%5 == 0 }),
		wnep(func(in int) bool { return in%7 == 0 }),
	)
	checkPredicate(t, f2, 70, true)
	checkPredicate(t, f2, 10, false)
	checkPredicate(t, f2, 2, false)
	checkPredicate(t, f2, 3, false)
	checkPredicate(t, f2, 5, false)
	checkPredicate(t, f2, 7, false)

	f3 := and(
		wnep(func(s string) bool { return strings.Contains("abcdef", s) }),
		wnep(func(s string) bool { return strings.Contains("abcde", s) }),
		wnep(func(s string) bool { return strings.Contains("abcd", s) }),
		wnep(func(s string) bool { return strings.Contains("abc", s) }),
	)
	checkPredicate(t, f3, "a", true)
	checkPredicate(t, f3, "ab", true)
	checkPredicate(t, f3, "abc", true)
	checkPredicate(t, f3, "abcd", false)
	checkPredicate(t, f3, "abcde", false)
	checkPredicate(t, f3, "abcdef", false)
	checkPredicate(t, f3, "abcdefg", false)
	checkPredicate(t, f3, "g", false)
	checkPredicate(t, f3, "", true)

	f4 := and((func(int) (bool, error))(nil), (func(int) (bool, error))(nil))
	if f4 != nil {
		t.Fatalf("must be nil")
	}
}

func TestOr(t *testing.T) {
	f := or(
		wnep(func(in int) bool { return in%2 == 0 }),
		wnep(func(in int) bool { return in%5 == 0 }),
	)
	checkPredicate(t, f, 10, true)
	checkPredicate(t, f, 2, true)
	checkPredicate(t, f, 5, true)
	checkPredicate(t, f, 3, false)

	f2 := or(
		wnep(func(in int) bool { return in%2 == 0 }),
		wnep(func(in int) bool { return in%5 == 0 }),
		wnep(func(in int) bool { return in%7 == 0 }),
	)
	checkPredicate(t, f2, 70, true)
	checkPredicate(t, f2, 10, true)
	checkPredicate(t, f2, 2, true)
	checkPredicate(t, f2, 5, true)
	checkPredicate(t, f2, 7, true)
	checkPredicate(t, f2, 3, false)

	f3 := or(
		wnep(func(s string) bool { return strings.Contains("abcdef", s) }),
		wnep(func(s string) bool { return strings.Contains("abcde", s) }),
		wnep(func(s string) bool { return strings.Contains("abcd", s) }),
		wnep(func(s string) bool { return strings.Contains("abc", s) }),
	)
	checkPredicate(t, f3, "a", true)
	checkPredicate(t, f3, "ab", true)
	checkPredicate(t, f3, "abc", true)
	checkPredicate(t, f3, "abcd", true)
	checkPredicate(t, f3, "abcde", true)
	checkPredicate(t, f3, "abcdef", true)
	checkPredicate(t, f3, "abcdefg", false)
	checkPredicate(t, f3, "g", false)
	checkPredicate(t, f3, "", true)

	f4 := or((func(int) (bool, error))(nil), (func(int) (bool, error))(nil))
	if f4 != nil {
		t.Fatalf("must be nil.")
	}
}

func TestNotAndComparableEquals(t *testing.T) {
	f1 := Not(ComparableEquals(1))
	checkPredicate(t, f1, 1, false)
	checkPredicate(t, f1, 2, true)

	f2 := Not((func(int) (bool, error))(nil))
	if f2 != nil {
		t.Error("f2 must be nil.")
	}
}

func TestNewPredicates(t *testing.T) {
	f := func(int) (bool, error) { return true, nil }
	f1 := newPredicates(f, f)
	if f1 == nil {
		t.Error("must not be nil")
	}

	f2 := newPredicates(f, f, f)
	if f2 == nil {
		t.Error("must not be nil")
	}

	f3 := newPredicates(f, f, f, f)
	if f3 == nil {
		t.Error("must not be nil")
	}

	fn := (func(int) (bool, error))(nil)
	f4 := newPredicates(fn, f)
	if f4 != nil {
		t.Error("must be nil")
	}

	f5 := newPredicates(f, fn)
	if f5 != nil {
		t.Error("must be nil")
	}

	f6 := newPredicates(f, f, fn)
	if f6 != nil {
		t.Error("must be nil")
	}

}
