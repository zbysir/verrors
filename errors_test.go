package verrors

import (
	"errors"
	"fmt"
	"testing"
)

func TestVerrors(t *testing.T) {
	root := errors.New("mysql cannot connect")
	err := WithCode(root, 500)
	err = WithCode(err, 300) // the late code 300 will override 500

	err = fmt.Errorf("插入错误: %w", WithStack(WithStack(WithStack(err))))

	err = fmt.Errorf("请求错误: %w", WithCode(err, 600))

	fe := formatInternalError{err}

	t.Logf("detail : \n%+v", fe)

	if Cause(err) != root {
		t.Fatal("Cause is't root")
	}

	if err.Error() != "请求错误: 插入错误: mysql cannot connect" {
		t.Fatalf("wrong err.Error() : %s", err.Error())
	}

	if !HasCode(err, 300) {
		t.Fatal("err has't code 300")
	}

	if HasCode(err, 500) {
		t.Fatal("err has code 500")
	}

	if len(Unpack(err)) != 3 {
		t.Fatalf("wrong len for unpacked error :%v", len(Unpack(err)))
	}

	t.Logf("stack: \n%v", Stack(err))
	t.Logf("stackDeep: \n%v", StackDeep(err))

	t.Logf("hasCode 300: %v, hasCode 500: %v", HasCode(err, 300), HasCode(err, 500))
	t.Logf("msg: %s", err.Error())
	t.Logf("detail2: \n%+v", NewFormatError2(err))
}

func Test113(t *testing.T) {
	root := errors.New("mysql cannot connent")
	e1 := fmt.Errorf("插入错误 %w", root)
	e2 := fmt.Errorf("请求错误 %w", e1)

	t.Logf("%+v", e2)
	t.Logf("%v", e2)
}

// test in readme.md
func TestReadMe(t *testing.T) {
	{
		// NewError
		err := WithCode(errors.New("file not found"), 500)

		err = fmt.Errorf("check health error: %w", err)
		t.Log(StdPackErrorsFormatter(Unpack(err)))
	}

	{
		// WithCode
		err := errors.New("file not found")
		err = fmt.Errorf("check health error: %w", WithCode(err, 400))

		t.Log(StdPackErrorsFormatter(Unpack(err)))
	}

	{
		// WithValue
		err := errors.New("file not found")
		err = fmt.Errorf("check health error: %w", WithValue(err, "retry", true))

		t.Log(StdPackErrorsFormatter(Unpack(err)))
	}

	{
		// formatPackErrors2
		StdPackErrorsFormatter = formatPackErrors2

		err := errors.New("file not found")
		err = fmt.Errorf("check health error: %w", WithStack(WithCode(err, 400)))

		t.Logf("\n%+v", WithFormat(err))
	}

	{
		// errors.Unwrap()
		root := errors.New("file not found")
		err := ToInternalError(fmt.Errorf("check health error: %w", WithCode(WithStack(root), 400)))
		t.Log(errors.Unwrap(err) == root) // true
	}
}

func TestStdUnwrap(t *testing.T) {
	{
		root := errors.New("file not found")
		err := WithStack(fmt.Errorf("check health error: %w", root))
		if errors.Unwrap(err) != root {
			t.Fatalf("Unwrap error, want root, but %T", errors.Unwrap(err))
		}
	}
	{
		root := errors.New("file not found")
		err := WithStack(fmt.Errorf("check health error: %v", WithStack(root)))
		if errors.Unwrap(err) != nil {
			t.Fatalf("Unwrap error, want nil, but %T", errors.Unwrap(err))
		}
	}
}
