package verrors

import (
	"errors"
	"testing"
)

func TestErrorf(t *testing.T) {
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

	{
		err := errors.New("file not found")

		err = Errorfc(500, "check health error: %w", err)
		err = Errorf("doSomeThingAService err: %w", err)

		t.Logf("Errorfc & Errorf:\n%+v", err)

		// - doSomeThingAService err: check health error: file not found [ stack = github.com/zbysir/verrors.TestErrorfc Z:/go_project/verrors/extra_test.go:29 ]
		// - check health error: file not found [ code = 500; stack = github.com/zbysir/verrors.TestErrorfc Z:/go_project/verrors/extra_test.go:28 ]
		// - file not found
	}
}
