package verrors

import "fmt"

type Wrapper interface {
	Unwrap() error
}

// InternalError 将不会在unpack里识别为独立一层错误, 而是作为上层错误的子数据.
type InternalError interface {
	// InternalError和Unwrap的行为不同
	// 在InternalError中实现的Unwrap一般都会跳过当层错误, 为了实现InternalError对Unwrap透明.
	// InternalError一般返回真正的下一层错误, 不过也有例外,
	//   如toInternalError就是一个对Unwrap和InternalError都透明处理的错误, 目的是实现将任何的错误都包裹为InternalError.
	// 你可以在看到Errorf()方法中看到WrapInterError的用例.
	// InternalError只会在Unpack中使用, 而Unwrap会做为通用解包方法.
	InternalError() error
}

var _ InternalError = toInternalError{}

type toInternalError struct {
	err error
}

func (e toInternalError) Unwrap() error {
	return Unwrap(e.err)
}

func (e toInternalError) InternalError() error {
	if u, ok := e.err.(Wrapper); ok {
		return u.Unwrap()
	}
	return nil
}

func (e toInternalError) Error() string {
	return e.err.Error()
}

func (e toInternalError) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			_, _ = f.Write([]byte(StdPackErrorsFormatter(Unpack(e))))
			return
		}
	}
	_, _ = f.Write([]byte(e.Error()))
}

// ToInternalError return a error that implement InternalError and Wrapper interface.
// toInternalError用于将任何一个错误包裹成为InternalError.
//
// 如下代码:
//
// err := WithCode(errors.New("not found"), 400)
//
// 在Unpack的时候, 由于errors.New返回的错误没有实现InternalError, 所以始终会被当做两层错误, 但实则用户只是想使用WithCode为error添加一个code码信息, 而不是Wrap,
// 使用toInternalError就能将errors.New返回的错误包裹成为一个InternalError, 同时还能兼容到fmt.Errorf("%w")
//
// 可以移步到 extra_test.go 中查看测试用例以理解它
func ToInternalError(err error) error {
	return toInternalError{err}
}

// Unwrap会解开一个err, 并且跳过当中的internalError, 返回底层错误.
// 用在所有的InternalError.Unwrap中, 用来兼容官方errors.Unwrap.
func Unwrap(err error) (next error) {
	// 如果err是一个internalError, 则不用unwrap, 而是直接放入internal并继续解析接下来的internalError
	// 否则就unwrap, 再继续解析接下来的internalError
	_, isInternal := err.(InternalError)
	if !isInternal {
		if w, ok := err.(Wrapper); ok {
			err = w.Unwrap()
		} else {
			return nil
		}
	}

	// 如果自己是InternalError, 则返回的internal包含自己
	for i, ok := err.(InternalError); ok; i, ok = err.(InternalError) {
		err = i.InternalError()
	}

	// 如果err是internalError, 说明还没真正的Unwrap过, 需要再Unwrap
	// 这是为了处理官方errors.Unwrap兼容的问题:
	// root:= errors.New("not found")
	// errors.Unwrap(WithStack(fmt.Errorf("xx :%w", root))) 应该返回 root
	if isInternal {
		if err != nil {
			if w, ok := err.(Wrapper); ok {
				err = w.Unwrap()
			} else {
				err = nil
			}
		}
	}

	return err
}
