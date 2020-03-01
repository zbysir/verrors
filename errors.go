package verrors

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

// 跳过下一层.
func (w toInternalError) Unwrap() error {
	if u, ok := w.err.(Wrapper); ok {
		return u.Unwrap()
	}
	return nil
}

func (w toInternalError) InternalError() error {
	return w.Unwrap()
}

func (w toInternalError) Error() string {
	return w.err.Error()
}

// NewToInternalError return a error that implement InternalError and Wrapper interface.
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
func NewToInternalError(err error) error {
	return toInternalError{err}
}

// Unwrap会解开一个err, 并且跳过当中的internalError, 返回底层错误.
func Unwrap(err error) (next error, internal []error) {
	// 如果err是一个internalError, 则不用unwrap, 而是直接放入internal并继续解析接下来的internalError
	// 否则就unwrap, 再继续解析接下来的internalError
	_, isInternal := err.(InternalError)
	if !isInternal {
		if w, ok := err.(Wrapper); ok {
			err = w.Unwrap()
		} else {
			return nil, nil
		}
	}

	// 如果自己是InternalError, 则返回的internal包含自己
	for i, ok := err.(InternalError); ok; i, ok = err.(InternalError) {
		internal = append(internal, err)
		err = i.InternalError()
	}

	return err, internal
}
