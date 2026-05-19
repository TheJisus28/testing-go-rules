package errors

import "fmt"

const (
	CodeNotFound            = "not_found"
	CodeValidation          = "validation_error"
	CodeUnauthorized        = "unauthorized"
	CodeForbidden           = "forbidden"
	CodeConflict            = "conflict"
	CodeAlreadyExists       = "already_exists"
	CodeUnprocessable       = "unprocessable_entity"
	CodeInternal            = "internal_error"
	StatusSuccess           = "success"
	StatusReasonOK          = "ok"
)

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func Validation(msg string) *DomainError {
	return &DomainError{Code: CodeValidation, Message: msg}
}

func NotFound(msg string) *DomainError {
	return &DomainError{Code: CodeNotFound, Message: msg}
}

func Unauthorized(msg string) *DomainError {
	return &DomainError{Code: CodeUnauthorized, Message: msg}
}

func Forbidden(msg string) *DomainError {
	return &DomainError{Code: CodeForbidden, Message: msg}
}

func Conflict(msg string) *DomainError {
	return &DomainError{Code: CodeConflict, Message: msg}
}

func AlreadyExists(msg string) *DomainError {
	return &DomainError{Code: CodeAlreadyExists, Message: msg}
}

func Unprocessable(msg string) *DomainError {
	return &DomainError{Code: CodeUnprocessable, Message: msg}
}

func Internal(msg string) *DomainError {
	return &DomainError{Code: CodeInternal, Message: msg}
}
