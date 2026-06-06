package common

type ViolationType string

const NotEmptyViolationType ViolationType = "not_empty"
const EmailViolationType ViolationType = "email"
const ConflictViolationType ViolationType = "conflict"
const InvalidFormatViolationType ViolationType = "invalid_format"

type FieldError struct {
	Key   string
	Error string
	Value any
}

type ValidationError struct {
	FieldErrors []FieldError
}

func (err ValidationError) Error() string {
	return "invalid data"
}
