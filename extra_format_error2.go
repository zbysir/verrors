package verrors

import (
	"fmt"
	"strings"
)

// formatInternalError implement fmt.Formatter(use %+v) to print formatted error info.
// 打印的log格式如下:
// - [600] 请求错误: 插入错误: mysql cannot connect >> go.zhuzi.me/go/errors/verrors.TestVerrors Z:/golang/go_path/src/go.zhuzi.me/go/errors/verrors/errors_test.go:17
// - [300] 插入错误: mysql cannot connect >> go.zhuzi.me/go/errors/verrors.TestVerrors Z:/golang/go_path/src/go.zhuzi.me/go/errors/verrors/errors_test.go:15
// - mysql cannot connect
type formatInternalError2 struct {
	err error
}

// 跳过下一层.
func (e formatInternalError2) Unwrap() error {
	return Unwrap(e.err)
}

func (e formatInternalError2) InternalError() error {
	return e.Unwrap()
}

func (e formatInternalError2) Error() string {
	return e.err.Error()
}

// 简单的打印错误, 只是为了方便临时查看, 建议用户实现自己的formatInternalError打印方法.
// use %+v to print more info.
func (e formatInternalError2) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			_, _ = f.Write([]byte(formatPackErrors2(Unpack(e))))
			return
		}
	}
	_, _ = f.Write([]byte(e.Error()))
}

func NewFormatError2(err error) error {
	return formatInternalError2{err}
}

// formatPackErrors2 返回的文本 包括错误消息 / code / stack
func formatPackErrors2(ps PackErrors) string {
	var s strings.Builder
	for _, v := range ps {
		if s.Len() != 0 {
			s.WriteString("\n")
		}
		s.WriteString("- ")

		code, codeExist := v.Get("code")
		if codeExist {
			s.WriteString(fmt.Sprintf("[%v] ", code))
		}

		s.WriteString(fmt.Sprintf("%v", v.Cause()))

		loc, locExist := v.Get("stack")
		if locExist {
			s.WriteString(fmt.Sprintf(" >> %s", loc))
		}
	}

	return s.String()
}
