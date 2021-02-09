package core

import (
	"context"
	"fmt"
)

type validationError struct {
	error
}

func AsValidationErr(err error) error {
	return &validationError{
		error: err,
	}
}

func IsValidationError(err error) bool {
	_, ok := err.(*validationError)
	return ok
}

type ContextAwareError struct {
	Wrapped error
	Ctx     context.Context
}

func (cae *ContextAwareError) Error() string {
	return cae.Wrapped.Error()
}

func CtxAware(ctx context.Context, err error) error {
	return &ContextAwareError{
		Wrapped: err,
		Ctx:     ctx,
	}
}

type errNotFound struct {
	error
}

func (nf *errNotFound) Error() string {
	return fmt.Sprintf("not found: %v", nf)
}

func ErrNotFound(err error) error {
	return &errNotFound{
		err,
	}
}

func IsNotFound(err error) bool {
	_, ok := err.(*errNotFound)
	return ok
}
