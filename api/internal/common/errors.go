package common

import "fmt"

type ErrUnauthorised struct {
	Message string
}

func (err ErrUnauthorised) Error() string {
	if err.Message != "" {
		return err.Message
	}

	return "the current user is not authorized to perform this action"
}

type ErrNotFound struct {
	Message string
}

func (err ErrNotFound) Error() string {
	if err.Message != "" {
		return err.Message
	}

	return "the requested data does not exist"
}

// ErrConflict is the base type for any unique-constraint violation.
type ErrConflict struct {
	Message string
}

func (e ErrConflict) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "a conflict occurred"
}

// ErrEnvironmentNameAlreadyInUse wraps ConflictError to add specificity.
type ErrEnvironmentNameAlreadyInUse struct {
	Name  string
	cause ErrConflict
}

func (e ErrEnvironmentNameAlreadyInUse) Error() string {
	return fmt.Sprintf("environment name %q is already in use", e.Name)
}

// Unwrap lets errors.Is/errors.As walk up to ConflictError.
func (e ErrEnvironmentNameAlreadyInUse) Unwrap() error {
	return e.cause
}

// Constructor — callers never build the struct directly.
func NewErrEnvironmentNameAlreadyInUse(name string) ErrEnvironmentNameAlreadyInUse {
	return ErrEnvironmentNameAlreadyInUse{
		Name:  name,
		cause: ErrConflict{Message: "environment name already in use"},
	}
}
