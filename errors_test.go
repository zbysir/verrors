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

func TestStdUnwrap(t *testing.T) {
	root := errors.New("mysql cannot connect")
	err := WithCode(WithStack(root), 500)

	var unwraped []error
	for err != nil {
		unwraped = append(unwraped, err)
		t.Logf("- %+v", err)

		err = errors.Unwrap(err)
	}

	t.Logf("has %v error in the chain", len(unwraped))
}

// test for readme.md
func TestReadMe(t *testing.T) {
	err := WithCode(errors.New("file not found"), 500)

	err = fmt.Errorf("check health error: %w", err)
	t.Log(StdPackErrorsFormatter(Unpack(err)))

	err = errors.New("file not found")
	err = fmt.Errorf("check health error: %w", WithCode(err, 400))

	t.Log(StdPackErrorsFormatter(Unpack(err)))

	err = errors.New("file not found")
	err = fmt.Errorf("check health error: %w", WithValue(err, "retry", true))

	t.Log(StdPackErrorsFormatter(Unpack(err)))

	StdPackErrorsFormatter = formatPackErrors2

	err = errors.New("file not found")
	err = fmt.Errorf("check health error: %w", WithStack(WithCode(err, 400)))

	t.Logf("\n%+v", NewFormatError(NewToInternalError(err)))
}
