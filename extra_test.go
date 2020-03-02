package verrors

import (
	"errors"
	"testing"
)

func TestErrorfc(t *testing.T) {
	err := errors.New("file not found")
	testFileName := "/usr/xx.txt"

	ev := Errorfc(400, "do something err: %v, fileName: %s", err, testFileName)
	t.Logf("format with v:\n%+v", ev)

	// only has on error in chain:
	// - do something err: file not found, fileName: /usr/xx.txt [ code = 400; stack = github.com/zbysir/verrors.TestErrorfc Z:/go_project/verrors/extra_test.go:12 ]

	ew := Errorfc(400, "do something err: %w, fileName: %s", err, testFileName)
	t.Logf("format with w:\n%+v", ew)

	// has two error in chain:
	// - do something err: file not found, fileName: /usr/xx.txt [ code = 400; stack = github.com/zbysir/verrors.TestErrorfc Z:/go_project/verrors/extra_test.go:18 ]
	// - file not found
}

func TestErrorf(t *testing.T) {
	err := errors.New("file not found")

	e := Errorf("do something err: %w", err)

	t.Logf("\n%+v", e)
}
