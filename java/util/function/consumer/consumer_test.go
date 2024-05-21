package consumer

import (
	"bytes"
	"testing"
)

func TestWrapNoErr(t *testing.T) {
	f1 := WrapNoErr[any](nil)
	if f1 != nil {
		t.Error("should be nil.")
	}
	f2 := WrapNoErr(func(i int) {
		t.Helper()
		if i != 1 {
			t.Errorf("want=1, got=%d", i)
		}
	})
	if f2 == nil {
		t.Fatalf("must not be nil")
	}
	err := f2(1)
	if err != nil {
		t.Errorf("must be nil")
	}
}

func TestCompose(t *testing.T) {

	type data struct {
		first, second, third, forth bool
	}
	first := func(d *data) error {
		d.first = true
		return nil
	}
	second := func(d *data) error {
		d.second = true
		return nil
	}
	third := func(d *data) error {
		d.third = true
		return nil
	}
	forth := func(d *data) error {
		d.forth = true
		return nil
	}

	check := func(t *testing.T, c Consumer[*data], want data) {
		t.Helper()
		in := data{}
		err := c(&in)
		if err != nil {
			t.Errorf("must not return error: %s", err)
			return
		}
		if in != want {
			t.Errorf("want=%+v, got=%+v", want, in)

		}
	}
	check(t, Compose(first), data{first: true})
	check(t, Compose(first, second), data{first: true, second: true})
	check(t, Compose(first, second, third), data{first: true, second: true, third: true})
	check(t, Compose(first, second, third, forth), data{first: true, second: true, third: true, forth: true})
}

type adjustChecker[T comparable] struct {
	t    *testing.T
	want T
}

func (c *adjustChecker[T]) do(in T) error {
	c.t.Helper()
	if in != c.want {
		c.t.Errorf("want=%#v, got=%#v", c.want, in)
	}
	return nil
}

func checkAdjust[T any, U comparable](t *testing.T, in T, want U) {
	t.Helper()
	cc := &adjustChecker[U]{t, want}
	f := Adjust[T, U](cc.do)
	if f == nil {
		t.Fatalf("should be nil.")

	}
	f(in)
}

type (
	myString string
	myInt    int
)

func TestCast(t *testing.T) {
	checkAdjust(t, myString("mystring"), "mystring")
	checkAdjust(t, "mystring", myString("mystring"))
	checkAdjust(t, myInt(1), 1)
	checkAdjust(t, 1, myInt(1))

	f1 := Adjust[any, any](nil)
	if f1 != nil {
		t.Error("must be nil")
	}
	f2 := Adjust[int, *bytes.Buffer](func(*bytes.Buffer) error { return nil })
	if f2 != nil {
		t.Error("must be nil")
		return
	}
}
