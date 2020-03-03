package verrors

import "errors"

// WithCode do the same thing as `WithValue(err, "code", code)`
func WithCode(err error, code int) error {
	return WithValue(err, "code", code)
}

func NewCode(msg string, code int) error {
	return WithCode(errors.New(msg), code)
}
