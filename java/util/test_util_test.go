package util

import "testing"

type called struct {
	t        *testing.T
	isCalled bool
	err      error
}

func (c *called) f() error {
	c.isCalled = true
	return c.err
}

func (c *called) checkCalled() {
	c.t.Helper()
	if !c.isCalled {
		c.t.Error("not called")
	}
}

func (c *called) checkNotCalled() {
	c.t.Helper()
	if c.isCalled {
		c.t.Error("should not called")
	}
}

type consumerCalled[T comparable] struct {
	called
	want T
}

func newConsumerCalled[T comparable](t *testing.T, want T, err error) *consumerCalled[T] {
	return &consumerCalled[T]{called{t, false, err}, want}
}

func (c *consumerCalled[T]) consume(got T) error {
	c.t.Helper()
	if got != c.want {
		c.t.Errorf("want=%v, got=%v", c.want, got)
	}
	return c.f()
}

type runnableCalled struct {
	called
}

func newRunnableCalled(t *testing.T, err error) *runnableCalled {
	return &runnableCalled{called{t, false, err}}
}

func (r *runnableCalled) run() error {
	r.t.Helper()
	return r.f()
}
