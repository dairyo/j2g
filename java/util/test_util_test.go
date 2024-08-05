package util

import "testing"

type called struct {
	t      *testing.T
	called bool
	err    error
}

func (c *called) f() error {
	c.called = true
	return c.err
}

func (c *called) checkCalled() {
	c.t.Helper()
	if !c.called {
		c.t.Error("not called")
	}
}

func (c *called) checkNotCalled() {
	c.t.Helper()
	if c.called {
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
	c.called.t.Helper()
	if got != c.want {
		c.t.Errorf("want=%v, got=%v", c.want, got)
	}
	return c.f()
}
