package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zbysir/verrors"
	"testing"
)

type withStatusError struct {
	err    error
	status int64
}

// 实现wrap接口可以兼容官方包方法, 如errors.Is
func (w withStatusError) Unwrap() error {
	return w.err
}

// 它是一个错误
func (w withStatusError) Error() string {
	return w.err.Error()
}

// 实现InternalError表明是内部错误, 不会被当为一层
func (w withStatusError) InternalError() {
}

// 在unpack的时候会被当做数据赋值给本层的错误
func (w withStatusError) Set(s verrors.Store) {
	s.Set("status", w.status)
}

func TestStatusError(t *testing.T) {
	e := errors.New("mysql cannot connect")
	e = fmt.Errorf("数据库错误: %w", e)
	e = fmt.Errorf("插入错误: %w", verrors.WithStack(e))
	e = fmt.Errorf("系统错误: %w", withStatusError{err: e, status: 300})

	up := verrors.Unpack(e)
	bs, _ := json.MarshalIndent(up, " ", " ")
	t.Logf("%s", bs)
}
