package function

import (
	"errors"
	"strconv"
	"testing"
)

func f[T any, U any](f func(T) (U, error)) Function[T, U] {
	return NewFunctionFunc(f)
}

func fe[T any, U any](f func(T) U) Function[T, U] {
	return NewFunctionFuncWithNoError(f)
}

func c[T any, U any, V any](f1 Function[T, U], f2 Function[U, V]) Function[T, V] {
	return Concatenate(f1, f2)
}

func check[T any, U comparable](t *testing.T, f Function[T, U], in T, want U) {
	t.Helper()
	got, err := f.Apply(in)
	if err != nil {
		t.Fatalf("must not return error but %q", err)
	}
	if want != got {
		t.Fatalf("want=%v, got=%v", want, got)
	}
}

func checkError[T any, U comparable](t *testing.T, f Function[T, U], in T, want error) {
	t.Helper()
	v, got := f.Apply(in)
	if v != *new(U) {
		t.Fatalf("must not return value but %v", v)
	}
	if got != want {
		t.Fatalf("error mismatch: want=%q, got=%q", want, got)
	}
}

type eq struct {
	i int
}

func (e *eq) Equals(i int) bool {
	return e.i == i
}

func TestNewFunctionFunc(t *testing.T) {
	// java's t -> t + 1
	f1 := f(func(i int) (int, error) { return i + 1, nil })
	check(t, f1, 1, 2)

	// java's static method usecase like Integer::parseInt.
	f2 := f(strconv.Atoi)
	check(t, f2, "100", 100)

	// java's instance method usecase like obj::Equals.
	v := &eq{1}
	f3 := f(func(i int) (bool, error) { return v.Equals(i), nil })
	check(t, f3, 1, true)
	check(t, f3, 2, false)
}

func TestNewFunctionFuncWithNoError(t *testing.T) {
	// java's i -> i + 1
	f1 := fe(func(i int) int { return i + 1 })
	check(t, f1, 1, 2)

	// java's static method usecase like Integer::parseInt.
	f2 := fe(strconv.Itoa)
	check(t, f2, 100, "100")

	// java's instance method usecase like obj::Equals.
	v := &eq{1}
	f3 := fe(func(i int) bool { return v.Equals(i) })
	check(t, f3, 1, true)
	check(t, f3, 2, false)
}

func checkNil[T any, U any](t *testing.T, f Function[T, U]) {
	t.Helper()
	if f != nil {
		t.Error("f must be nil")
	}
}

func TestConcatenate(t *testing.T) {
	// java's f1.andThen(f2) or f1.compose(f2).
	f1 := c(fe(strconv.Itoa), f(strconv.Atoi))
	check(t, f1, 100, 100)

	// java's f1.andThen(f2).andThen(f3).
	f2 := c(fe(strconv.Itoa), c(f(strconv.Atoi), fe(strconv.Itoa)))
	check(t, f2, 100, "100")
	f3 := c(c(fe(strconv.Itoa), f(strconv.Atoi)), fe(strconv.Itoa))
	check(t, f3, 100, "100")

	checkNil(t, c(Function[int, string](nil), f(strconv.Atoi)))
	checkNil(t, c(fe(strconv.Itoa), Function[string, int](nil)))
	checkNil(t, c(fe(strconv.Itoa), c(Function[string, int](nil), fe(strconv.Itoa))))

	want := errors.New("foo")
	f4 := c(f(strconv.Atoi), f(func(int) (int, error) { return 0, want }))
	checkError(t, f4, "100", want)
	f5 := c(f(func(int) (int, error) { return 0, want }), fe(strconv.Itoa))
	checkError(t, f5, 100, want)
}

func TestIdentity(t *testing.T) {
	f := Identity[int]()
	check(t, f, 100, 100)
}
