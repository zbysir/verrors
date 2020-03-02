package verrors

import "fmt"

func HasCode(err error, code int) bool {
	for err != nil {
		var pg PackError
		err, pg = UnpackOnce(err)
		c, _ := pg.Get("code")
		if c == code {
			return true
		}
	}

	return false
}

// Stack 返回最近错误发生的位置
func Stack(err error) string {
	for err != nil {
		var pg PackError
		err, pg = UnpackOnce(err)
		loc, locExist := pg.Get("stack")
		if locExist {
			return loc.(string)
		}
	}

	return ""
}

// StackDeep 返回最底层错误发生的位置
func StackDeep(err error) string {
	ps := Unpack(err)
	loc := ""

	for _, v := range ps {
		l, _ := v.Get("stack")
		if l != nil && l != "" {
			loc = l.(string)
		}
	}

	return loc
}

// 最底层错误
func Cause(err error) error {
	for {
		e := Unwrap(err)
		if e != nil {
			err = e
		} else {
			break
		}
	}

	return err
}

// Errofc is shorthand for WithStack/WithCode/fmt.Errorf
func Errorfc(code int, format string, args ...interface{}) (r error) {
	return WithStack(WithCode(ToInternalError(fmt.Errorf(format, args...)), code), 2)
}

func Errorf(format string, args ...interface{}) error {
	return ToInternalError(fmt.Errorf(format, args...))
}
