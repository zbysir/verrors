package verrors

import "fmt"

// formatInternalError implement fmt.Formatter(use %+v) to print formatted error info.
// info include all value, like following log:
//   - get User error: mysql can't connect [ code = 400; stack = go.zhuzi.me/go/errors/verrors.TestFormat Z:/golang/go_path/src/go.zhuzi.me/go/errors/verrors/extra_format_error_test.go:11 ]
//   - mysql can't connect
type formatInternalError struct {
	err error
}

// 跳过下一层.
func (e formatInternalError) Unwrap() error {
	if u, ok := e.err.(Wrapper); ok {
		return u.Unwrap()
	}
	return nil
}

func (e formatInternalError) InternalError() error {
	return e.err
}

func (e formatInternalError) Error() string {
	return e.err.Error()
}

// 简单的打印错误, 只是为了方便临时查看, 建议用户实现自己的formatInternalError打印方法.
// use %+v to print more info.
func (e formatInternalError) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			_, _ = f.Write([]byte(StdPackErrorsFormatter(Unpack(e))))
			return
		}
	}
	_, _ = f.Write([]byte(e.Error()))
}

type PackErrorsFormatter func(e PackErrors) string

// StdPackErrorsFormatter is standard PackErrors formatter, will print like following text:
//
//   - get User error: mysql can't connect [ code = 400; stack = go.zhuzi.me/go/errors/verrors.TestFormat Z:/golang/go_path/src/go.zhuzi.me/go/errors/verrors/extra_format_error_test.go:11 ]
//   - mysql can't connect
//
// 你可以覆盖这个值来实现覆盖verrors.Errorf()返回的错误的打印行为.
var StdPackErrorsFormatter PackErrorsFormatter

func NewFormatError(err error) error {
	return formatInternalError{err}
}

func init() {
	StdPackErrorsFormatter = formatPackErrors
}
