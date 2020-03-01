package verrors

import (
	"errors"
	"fmt"
	"testing"
)

func TestWithValue(t *testing.T) {
	root := errors.New("mysql cannot connent")
	err := fmt.Errorf("插入错误: %w", WithValue(root, "code", 400, "status", 500))
	t.Logf("%s", err)

	up := Unpack(err)
	t.Logf("%+v", up)
	t.Logf("%v", up)
	//t.Logf("Detail : \n%s", Detail(err))
}
