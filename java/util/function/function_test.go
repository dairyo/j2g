package function

import (
	"errors"
	"strconv"
	"testing"
)

func c[T any, U any, V any](f1 Function[T, U], f2 Function[U, V]) Function[T, V] {
	return Compose(f1, f2)
}

func ne[T any, U any](f func(in T) U) Function[T, U] {
	return WrapNoErrFunc(f)
}

func checkFunction[T any, U comparable](t *testing.T, f Function[T, U], in T, want U) {
	t.Helper()
	got, err := f(in)
	if err != nil {
		t.Fatalf("must not return error but %q", err)
	}
	if want != got {
		t.Fatalf("want=%v, got=%v", want, got)
	}
}

func checkFunctionError[T any, U comparable](t *testing.T, f Function[T, U], in T, want error) {
	t.Helper()
	v, got := f(in)
	if v != *new(U) {
		t.Fatalf("must not return value but %v", v)
	}
	if got != want {
		t.Fatalf("error mismatch: want=%q, got=%q", want, got)
	}
}

type calc struct {
	i int
}

func (c *calc) Add(i int) (int, error) {
	return c.i + i, nil
}

func (c *calc) Sub(i int) int {
	return c.i - i
}

func TestNewFunctionFunc(t *testing.T) {
	// java's t -> t + 1
	f1 := func(i int) (int, error) { return i + 1, nil }
	checkFunction(t, f1, 1, 2)

	// java's static method usecase like Integer::parseInt.
	f2 := strconv.Atoi
	checkFunction(t, f2, "100", 100)

	// java's instance method usecase like obj::Equals.
	c := &calc{1}
	f3 := c.Add
	checkFunction(t, f3, 1, 2)
	checkFunction(t, f3, 2, 3)
}

func TestNewFunctionFuncWithNoError(t *testing.T) {
	// java's i -> i + 1
	f1 := ne(func(i int) int { return i + 1 })
	checkFunction(t, f1, 1, 2)

	// java's static method usecase like Integer::parseInt.
	f2 := ne((strconv.Itoa))
	checkFunction(t, f2, 100, "100")

	// java's instance method usecase like obj::Equals.
	c := &calc{1}
	f3 := ne((c.Sub))
	checkFunction(t, f3, 1, 0)
	checkFunction(t, f3, 2, -1)
}

func checkFunctionNil[T any, U any](t *testing.T, f Function[T, U]) {
	t.Helper()
	if f != nil {
		t.Error("f must be nil")
	}
}

func TestCompose(t *testing.T) {
	// java's f1.andThen(f2) or f1.compose(f2).
	f1 := c(ne(strconv.Itoa), strconv.Atoi)
	checkFunction(t, f1, 100, 100)

	// java's f1.andThen(f2).andThen(f3).
	f2 := c(ne(strconv.Itoa), c(strconv.Atoi, ne(strconv.Itoa)))
	checkFunction(t, f2, 100, "100")
	f3 := c(c(ne(strconv.Itoa), strconv.Atoi), ne(strconv.Itoa))
	checkFunction(t, f3, 100, "100")

	checkFunctionNil(t, c(Function[int, string](nil), strconv.Atoi))
	checkFunctionNil(t, c(ne(strconv.Itoa), Function[string, int](nil)))
	checkFunctionNil(t, c(ne(strconv.Itoa), c(Function[string, int](nil), ne(strconv.Itoa))))

	want := errors.New("foo")
	f4 := c(strconv.Atoi, func(int) (int, error) { return 0, want })
	checkFunctionError(t, f4, "100", want)
	f5 := c(func(int) (int, error) { return 0, want }, ne(strconv.Itoa))
	checkFunctionError(t, f5, 100, want)
}

func TestIdentity(t *testing.T) {
	f := Identity[int]()
	checkFunction(t, f, 100, 100)
}
