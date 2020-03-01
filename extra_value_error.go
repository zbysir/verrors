package verrors

import "fmt"

// 检查接口实现
var _ InternalError = valueInternalError{}
var _ Wrapper = valueInternalError{}
var _ Setter = valueInternalError{}

type valueInternalError struct {
	keys   []string
	values []interface{}
	err    error
}

func (e valueInternalError) Set(s Store) {
	for i, k := range e.keys {
		s.Set(k, e.values[i])
	}
}

func (e valueInternalError) Unwrap() error {
	if u, ok := e.err.(Wrapper); ok {
		return u.Unwrap()
	}
	return nil
}

func (e valueInternalError) InternalError() error {
	return e.err
}

func (e valueInternalError) Error() string {
	return e.err.Error()
}

func (e valueInternalError) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			_, _ = f.Write([]byte(StdPackErrorsFormatter(Unpack(e))))
			return
		}
	}
	_, _ = f.Write([]byte(e.Error()))
}

func WithValue(err error, keyValPairs ...interface{}) error {
	if len(keyValPairs)%2 == 1 {
		panic("bad length of kvPairs")
	}

	var v valueInternalError
	v.err = err
	l := len(keyValPairs)
	for i := 0; i < l; i += 2 {
		v.keys = append(v.keys, keyValPairs[i].(string))
		v.values = append(v.values, keyValPairs[i+1])
	}

	return v
}
