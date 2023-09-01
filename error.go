package uploader

import "fmt"

const (
	CONFLICT       = "conflict"
	INTERNAL       = "internal"
	INVALID        = "invalid"
	NOTFOUND       = "not_found"
	NOTIMPLEMENTED = "not_implemented"
	UNAUTHORIZED   = "unauthorized"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("uploader error: code=%s message=%s", e.Code, e.Message)
}

func Errorf(code, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
