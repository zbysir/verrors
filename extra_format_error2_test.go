package verrors

import (
	"errors"
	"fmt"
	"testing"
)

func TestFormat2(t *testing.T) {
	err := errors.New("mysql can't connect")
	err = fmt.Errorf("get User error: %w", WithCode(err, 400))
	err = NewFormatError2(err)

	t.Logf("\n%+v", err)
}
