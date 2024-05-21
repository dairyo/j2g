package supplier

import "testing"

func TestWrapNoErr(t *testing.T) {
	if WrapNoErr[any](nil) != nil {
		t.Error("must be nil.")
	}
	f := WrapNoErr[int](func() int { return 0 })
	if f == nil {
		t.Error("must not be nil")
	}
	got, err := f()
	if err != nil {
		t.Fatalf("must not return error: %s", err)
	}
	if got != 0 {
		t.Errorf("want=0, got=%d", got)
	}
}
