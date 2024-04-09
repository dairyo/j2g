package function

import (
	"errors"
	"strconv"
	"testing"
)

func f[T any, U any](f func(T) (U, error)) Applier[T, U] {
	return NewApplierFunc(f)
}

func fe[T any, U any](f func(T) U) Applier[T, U] {
	return NewApplierFuncWithNoError(f)
}

func c[T any, U any, V any](f1 Applier[T, U], f2 Applier[U, V]) Applier[T, V] {
	return Compose(f1, f2)
}

func checkApplier[T any, U comparable](t *testing.T, f Applier[T, U], in T, want U) {
	t.Helper()
	got, err := f.Apply(in)
	if err != nil {
		t.Fatalf("must not return error but %q", err)
	}
	if want != got {
		t.Fatalf("want=%v, got=%v", want, got)
	}
}

func checkApplierError[T any, U comparable](t *testing.T, f Applier[T, U], in T, want error) {
	t.Helper()
	v, got := f.Apply(in)
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

func TestNewApplierFunc(t *testing.T) {
	// java's t -> t + 1
	f1 := f(func(i int) (int, error) { return i + 1, nil })
	checkApplier(t, f1, 1, 2)

	// java's static method usecase like Integer::parseInt.
	f2 := f(strconv.Atoi)
	checkApplier(t, f2, "100", 100)

	// java's instance method usecase like obj::Equals.
	c := &calc{1}
	f3 := f(c.Add)
	checkApplier(t, f3, 1, 2)
	checkApplier(t, f3, 2, 3)
}

func TestNewApplierFuncWithNoError(t *testing.T) {
	// java's i -> i + 1
	f1 := fe(func(i int) int { return i + 1 })
	checkApplier(t, f1, 1, 2)

	// java's static method usecase like Integer::parseInt.
	f2 := fe(strconv.Itoa)
	checkApplier(t, f2, 100, "100")

	// java's instance method usecase like obj::Equals.
	c := &calc{1}
	f3 := fe(c.Sub)
	checkApplier(t, f3, 1, 0)
	checkApplier(t, f3, 2, -1)
}

func checkApplierNil[T any, U any](t *testing.T, f Applier[T, U]) {
	t.Helper()
	if f != nil {
		t.Error("f must be nil")
	}
}

func TestCompose(t *testing.T) {
	// java's f1.andThen(f2) or f1.compose(f2).
	f1 := c(fe(strconv.Itoa), f(strconv.Atoi))
	checkApplier(t, f1, 100, 100)

	// java's f1.andThen(f2).andThen(f3).
	f2 := c(fe(strconv.Itoa), c(f(strconv.Atoi), fe(strconv.Itoa)))
	checkApplier(t, f2, 100, "100")
	f3 := c(c(fe(strconv.Itoa), f(strconv.Atoi)), fe(strconv.Itoa))
	checkApplier(t, f3, 100, "100")

	checkApplierNil(t, c(Applier[int, string](nil), f(strconv.Atoi)))
	checkApplierNil(t, c(fe(strconv.Itoa), Applier[string, int](nil)))
	checkApplierNil(t, c(fe(strconv.Itoa), c(Applier[string, int](nil), fe(strconv.Itoa))))

	want := errors.New("foo")
	f4 := c(f(strconv.Atoi), f(func(int) (int, error) { return 0, want }))
	checkApplierError(t, f4, "100", want)
	f5 := c(f(func(int) (int, error) { return 0, want }), fe(strconv.Itoa))
	checkApplierError(t, f5, 100, want)
}

func TestIdentity(t *testing.T) {
	f := Identity[int]()
	checkApplier(t, f, 100, 100)
}
