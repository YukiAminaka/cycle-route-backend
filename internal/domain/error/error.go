package domainerror

import "errors"

// センチネルエラー
var (
	ErrValidation   = errors.New("validation error")
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

type Error struct {
	message string
	kind    error // センチネルエラーを保持
}

func (e *Error) Error() string { return e.message }

// Unwrap によって errors.Is でセンチネルエラーと比較できる
func (e *Error) Unwrap() error { return e.kind }

func New(message string, kind error) *Error {
	return &Error{message: message, kind: kind}
}
