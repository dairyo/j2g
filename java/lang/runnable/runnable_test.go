package runnable

import "testing"

func TestWrapNoErr(t *testing.T) {
	if WrapNoErr(nil) != nil {
		t.Error("must be nil.")
	}
	if WrapNoErr(func() {}) == nil {
		t.Error("must not be nil")
	}
}
