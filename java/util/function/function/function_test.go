package function

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"testing"
)

func c[T any, U any, V any](f1 Function[T, U], f2 Function[U, V]) Function[T, V] {
	return Compose(f1, f2)
}

func wne[T any, U any](f func(in T) U) Function[T, U] {
	return WrapNoErr(f)
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

func TestFunction(t *testing.T) {
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

func TestWrapNoErr(t *testing.T) {
	// java's i -> i + 1
	f1 := wne(func(i int) int { return i + 1 })
	checkFunction(t, f1, 1, 2)

	// java's static method usecase like Integer::parseInt.
	f2 := wne((strconv.Itoa))
	checkFunction(t, f2, 100, "100")

	// java's instance method usecase like obj::Equals.
	c := &calc{1}
	f3 := wne((c.Sub))
	checkFunction(t, f3, 1, 0)
	checkFunction(t, f3, 2, -1)

	f4 := WrapNoErr((func(int) int)(nil))
	if f4 != nil {
		t.Error("f4 must be nil.")
	}
}

func checkFunctionNil[T any, U any](t *testing.T, f Function[T, U]) {
	t.Helper()
	if f != nil {
		t.Error("f must be nil")
	}
}

func TestCompose(t *testing.T) {
	// java's f1.andThen(f2) or f1.compose(f2).
	f1 := c(wne(strconv.Itoa), strconv.Atoi)
	checkFunction(t, f1, 100, 100)

	// java's f1.andThen(f2).andThen(f3).
	f2 := c(wne(strconv.Itoa), c(strconv.Atoi, wne(strconv.Itoa)))
	checkFunction(t, f2, 100, "100")
	f3 := c(c(wne(strconv.Itoa), strconv.Atoi), wne(strconv.Itoa))
	checkFunction(t, f3, 100, "100")

	checkFunctionNil(t, c(Function[int, string](nil), strconv.Atoi))
	checkFunctionNil(t, c(wne(strconv.Itoa), Function[string, int](nil)))
	checkFunctionNil(t, c(wne(strconv.Itoa), c(Function[string, int](nil), wne(strconv.Itoa))))

	want := errors.New("foo")
	f4 := c(strconv.Atoi, func(int) (int, error) { return 0, want })
	checkFunctionError(t, f4, "100", want)
	f5 := c(func(int) (int, error) { return 0, want }, wne(strconv.Itoa))
	checkFunctionError(t, f5, 100, want)
}

func TestIdentity(t *testing.T) {
	f := Identity[int]()
	checkFunction(t, f, 100, 100)
}

func TestAdjust(t *testing.T) {
	f1 := func(b *bytes.Buffer) (*bytes.Buffer, error) {
		return b, nil
	}
	f2 := func(w io.Writer) (io.Writer, error) {
		return w, nil
	}
	b1 := &bytes.Buffer{}
	c := Compose(f1, Adjust[*bytes.Buffer, io.Writer, *bytes.Buffer, io.Writer](f2))
	if c == nil {
		t.Fatalf("must not be nil.")
	}
	b2, err := c(b1)
	if err != nil {
		t.Fatalf("must not return err: %s", err)
	}
	if b1 != b2 {
		t.Fatal("must be same.")
	}
}
