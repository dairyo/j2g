package util

import "errors"

var (
	ErrNilConsumer     = errors.New("Consumer is nil")
	ErrNilRunnable     = errors.New("Runnable is nil")
	ErrMapNilFunction  = errors.New("Function on Map argument is nil")
	ErrMapNilOptinal   = errors.New("Optional on Map argument is nil")
	ErrInvalidUsed     = errors.New("invalid optional is used")
	ErrEmpty           = errors.New("empty optional")
	ErrNilPredicate    = errors.New("Predicate is nil")
	ErrPredicateFailed = errors.New("Predicate returns false")
	ErrPredicateErr    = errors.New("Predicate returns error")
	ErrNilSupplier     = errors.New("Supplier is nil")
	ErrSupplierErr     = errors.New("Supplier returns error")
	ErrNoValue         = errors.New("Method is called for no value Optional")
)
