package internal

import (
	"bytes"
	"io"
	"testing"
)

type (
	myString string
	myInt    int
)

func checkCast[T any, U comparable](t *testing.T, in T, want U) {
	t.Helper()
	cast := Cast[T, U]()
	if cast == nil {
		t.Error("Cast returns nil.")
	}
	got, err := cast(in)
	if err != nil {
		t.Fatalf("must not return err: %s", err)
	}
	if want != got {
		t.Errorf("want=%#v, got=%#v", want, got)
	}
}

func TestCast(t *testing.T) {
	checkCast(t, myString("mystring"), "mystring")
	checkCast(t, "mystring", myString("mystring"))
	checkCast(t, myInt(1), 1)
	checkCast(t, 1, myInt(1))

	func() {
		cast := Cast[*bytes.Buffer, io.Writer]()
		if cast == nil {
			t.Error("fail to create cast function.")
			return
		}
		got, err := cast(&bytes.Buffer{})
		if err != nil {
			t.Fatalf("must not return err: %s", err)
		}
		if got == nil {
			t.Error("fail to cast.")
			return
		}
	}()

	func() {
		cast := Cast[io.Writer, *bytes.Buffer]()
		if cast == nil {
			t.Error("fail to create cast function.")
			return
		}
		b := &bytes.Buffer{}
		got, err := cast(b)
		if err != nil {
			t.Fatalf("must not return err: %s", err)
		}
		if got != b {
			t.Error("fail to cast.")
			return
		}

	}()

	func() {
		cast := Cast[*bytes.Buffer, io.Closer]()
		if cast != nil {
			t.Error("must return nil.")
			return
		}
	}()

}
