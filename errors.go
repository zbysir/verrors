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

// Unwrap会解开一个err, 并且跳过当中的internalError, 返回底层错误.
// 用在所有的InternalError.Unwrap中, 用来兼容官方errors.Unwrap.
func Unwrap(err error) (next error) {
	if _, isInternal := err.(InternalError); isInternal {
		// 跳过Internal
		for i, ok := err.(InternalError); ok; i, ok = err.(InternalError) {
			err = i.InternalError()
		}
	}

	if err != nil {
		if w, ok := err.(Wrapper); ok {
			err = w.Unwrap()
		} else {
			err = nil
		}
	}
	return err
}
