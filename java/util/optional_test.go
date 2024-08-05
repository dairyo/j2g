package util

import (
	"errors"
	"strconv"
	"testing"

	"github.com/dairyo/j2g/java/util/function/function"
)

func checkNotEmpty[T any](t *testing.T, o *Optional[T]) {
	t.Helper()
	if err := o.Error(); err != nil {
		t.Fatalf("NewOptional Failed: %s", err)
	}
}

func checkGet[T comparable](t *testing.T, o *Optional[T], want T) {
	t.Helper()
	got, err := o.Get()
	if err != nil {
		t.Fatalf("Get returns error: %s", err)
	}
	if got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestMap(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		i := NewOptional(int(1))
		checkNotEmpty(t, i)
		o := Map(i, function.WrapNoErr(strconv.Itoa))
		checkNotEmpty(t, o)
		checkGet(t, o, "1")
	})

	t.Run("string to int", func(t *testing.T) {
		i := NewOptional("1")
		checkNotEmpty(t, i)
		o := Map(i, strconv.Atoi)
		checkNotEmpty(t, o)
		checkGet(t, o, 1)
	})

	t.Run("argument Optional is nil", func(t *testing.T) {
		i := (*Optional[string])(nil)
		o := Map(i, strconv.Atoi)
		got := o.Error()
		if !errors.Is(got, ErrMapNilOptinal) {
			t.Fatalf("want=%s, got=%s", ErrMapNilOptinal, got)
		}
	})

	t.Run("argument function is nil", func(t *testing.T) {
		i := NewOptional("1")
		o := Map(i, (func(string) (int, error))(nil))
		got := o.Error()
		if !errors.Is(got, ErrMapNilFunction) {
			t.Fatalf("want=%s, got=%s", ErrMapNilFunction, got)
		}
	})

	t.Run("errorOptional", func(t *testing.T) {
		i := NewOptional((*int)(nil))
		o := Map(i, func(*int) (int, error) { return 1, nil })
		got := o.Error()
		if !errors.Is(got, ErrEmpty) {
			t.Fatalf("want=%s, got=%s", ErrEmpty, got)
		}
		if got.Error() != "invalid optional is passed: empty optional" {
			t.Fatalf("want=%s, got=%s", "invalid optional is passed: empty optional", got.Error())
		}
	})

	t.Run("function returns error", func(t *testing.T) {
		i := NewOptional("1")
		o := Map(i, func(string) (int, error) { return 0, errors.New("foo") })
		got := o.Error()
		if got.Error() != "function returns error: foo" {
			t.Fatalf("want=%s, got=%s", "function returns error: foo", got.Error())
		}
	})
}

func TestFlatMap(t *testing.T) {
	atoi := func(s string) (*Optional[int], error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return NewOptional(i), nil
	}

	t.Run("string to int", func(t *testing.T) {
		i := NewOptional("1")
		checkNotEmpty(t, i)
		o := FlatMap(i, atoi)
		checkNotEmpty(t, o)
		checkGet(t, o, 1)
	})

	t.Run("argument Optional is nil", func(t *testing.T) {
		i := (*Optional[string])(nil)
		o := FlatMap(i, atoi)
		got := o.Error()
		if !errors.Is(got, ErrMapNilOptinal) {
			t.Fatalf("want=%s, got=%s", ErrMapNilOptinal, got)
		}
	})

}

func TestIfPresent(t *testing.T) {
	t.Run("not empty success", func(t *testing.T) {
		c := newConsumerCalled(t, 1, nil)
		i := NewOptional(1)
		err := i.IfPresent(c.consume)
		c.checkCalled()
		if err != nil {
			t.Errorf("should not return error but %q.", err)
		}
	})

	t.Run("no consumer", func(t *testing.T) {
		i := NewOptional(1)
		err := i.IfPresent((func(int) error)(nil))
		if err != ErrNilConsumer {
			t.Errorf("should return %q, but %q.", ErrNilConsumer, err)
		}
	})

	t.Run("not empty error", func(t *testing.T) {
		want := errors.New("foo")
		c := newConsumerCalled(t, 1, want)
		i := NewOptional(1)
		err := i.IfPresent(c.consume)
		c.checkCalled()
		if err != want {
			t.Errorf("should return foo but %q.", err)
		}
	})

	t.Run("empty", func(t *testing.T) {
		c := newConsumerCalled[*int](t, nil, nil)
		i := NewOptional[*int](nil)
		err := i.IfPresent(c.consume)
		c.checkNotCalled()
		if err != ErrNoValue {
			t.Errorf("should return %q but %q.", ErrNoValue, err)
		}
	})
}
