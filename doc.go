// Package verrors can add any value for standard error in Go 1.13, like context.WithValue().
//
// This package has less code:
// - function: Unpack(err)
// - interface: Set(Store)
// - interface: InternalError() error
// simple but extensible,
// you can extend any value to error or custom print via they.
//
// I want to write some document in English,
// but my English is not satisfactory, you can visit test code in `errors_test.go` or README.md to understand them.
//
// e.g.
//   err := errors.New("file not found")
//   testFileName := "/usr/xx.txt"
//   err = verrors.Errorfc(400, "do something err: %w, fileName: %s", err, testFileName)
//   fmt.Printf("\n%+v", err)
//
//   will print:
//   - do something err: file not found, fileName: /usr/xx.txt [ code = 400; stack = verrors.TestErrorfc Z:/go/errors/verrors/extra_test.go:18 ]
//   - file not found
// The verrors.Errorfc is a sample function:
//   func Errorfc(code int,format string, args ...interface{}) (r error) {
//	   return NewFormatError(WithStack(WithCode(NewToInternalError(fmt.Errorf(format, args...)), code), 2))
//   }
package verrors // import "github.com/zbysir/verrors"
