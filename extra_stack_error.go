package verrors

import (
	"fmt"
	"runtime"
)

// 检查接口实现
var _ InternalError = stackInternalError{}
var _ Wrapper = stackInternalError{}
var _ Setter = stackInternalError{}

// 实现InternalError方法, 将不会在Detail里单独打印一行
type stackInternalError struct {
	err   error
	frame frame
}

func (e stackInternalError) Unwrap() error {
	if u, ok := e.err.(Wrapper); ok {
		return u.Unwrap()
	}
	return nil
}

func (e stackInternalError) Set(setter Store) {
	setter.Set("stack", e.frame.Format())
}

func (e stackInternalError) InternalError() error {
	return e.err
}

func (e stackInternalError) Error() string {
	return e.err.Error()
}

func (e stackInternalError) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			_, _ = f.Write([]byte(StdPackErrorsFormatter(Unpack(e))))
			return
		}
	}
	_, _ = f.Write([]byte(e.Error()))
}

// 包裹位置信息
func WithStack(err error, skip ...int) error {
	s := 1
	if len(skip) != 0 {
		s = skip[0]
	}
	return &stackInternalError{
		err:   err,
		frame: caller(s),
	}
}

// copy for xerror

// A Frame contains part of a call stack.
type frame struct {
	// Make room for three PCs: the one we were asked for, what it called,
	// and possibly a PC for skipPleaseUseCallersFrames. See:
	// https://go.googlesource.com/go/+/032678e0fb/src/runtime/extern.go#169
	frames [3]uintptr
}

// caller returns a Frame that describes a where on the caller's stack.
// The argument skip is the number of frames to skip over.
// caller(0) returns the where for the caller of caller.
func caller(skip int) frame {
	var s frame
	runtime.Callers(skip+1, s.frames[:])
	return s
}

// location reports the file, line, and function of a where.
//
// The returned function may be "" even if file and line are not.
func (f frame) location() (function, file string, line int) {
	frames := runtime.CallersFrames(f.frames[:])
	if _, ok := frames.Next(); !ok {
		return "", "", 0
	}
	fr, ok := frames.Next()
	if !ok {
		return "", "", 0
	}
	return fr.Function, fr.File, fr.Line
}

func (f frame) Format() string {
	var s string
	function, file, line := f.location()
	if function != "" {
		s += fmt.Sprintf("%s ", function)
	}
	if file != "" {
		s += fmt.Sprintf("%s:%d", file, line)
	}

	return s
}
