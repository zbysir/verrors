package verrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestUnpack(t *testing.T) {
	root := errors.New("mysql cannot connect")
	e := WithCode(root, 500)
	e = WithCode(e, 300)

	e = fmt.Errorf("插入错误: %w", WithStack(e))

	e = fmt.Errorf("请求错误: %w", WithCode(e, 600))

	up := Unpack(e)
	bs, _ := json.MarshalIndent(up, " ", " ")
	t.Logf("%s", bs)

	t.Logf("%v, %+v, %s", up, up, up)

	bs, _ = json.MarshalIndent(up.Merge(), " ", " ")
	t.Logf("Merge: %s", bs)

	return
}
