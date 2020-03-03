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
	return Unwrap(e.err)
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

	// 记住, 一定需要返回指针
	// 如果不返回指针, 在编译时并不会有错误, 但由于valueInternalError结构体包含切片无法比较,
	//   所以当比较(==)WithValue返回的error值的时候会panic, 而指针不会有这样的问题
	return &v
}
