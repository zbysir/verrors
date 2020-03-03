package verrors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestUnpack(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		packType []string
	}{
		{
			name:     "internal->error->internal->internal->error",
			err:      WithStack(fmt.Errorf("check health error: %w", WithStack(WithCode(errors.New("file not found"), 400)))),
			packType: []string{"*fmt.wrapError", "*errors.errorString"},
		},
		{
			name:     "error->internal->internal->error",
			err:      fmt.Errorf("check health error: %w", WithStack(WithCode(errors.New("file not found"), 400))),
			packType: []string{"*fmt.wrapError", "*errors.errorString"},
		},
		{
			name:     "error->error",
			err:      fmt.Errorf("check health error: %w", errors.New("file not found")),
			packType: []string{"*fmt.wrapError", "*errors.errorString"},
		},
		{
			name:     "error->error->error",
			err:      fmt.Errorf("check health error: %w", fmt.Errorf("x %w", errors.New("file not found"))),
			packType: []string{"*fmt.wrapError", "*fmt.wrapError", "*errors.errorString"},
		},
		{
			name:     "error->error->error",
			err:      fmt.Errorf("check health error: %w", fmt.Errorf("x %w", errors.New("file not found"))),
			packType: []string{"*fmt.wrapError", "*fmt.wrapError", "*errors.errorString"},
		},
		{
			name:     "error",
			err:      fmt.Errorf("check health error: %v", fmt.Errorf("x %w", errors.New("file not found"))),
			packType: []string{ "*errors.errorString"},
		},
		{
			name:     "internal->error->error",
			err:      WithStack(fmt.Errorf("x %w", errors.New("file not found"))),
			packType: []string{"*fmt.wrapError", "*errors.errorString"},
		},
		{
			name:     "internal->error",
			err:      WithStack(fmt.Errorf("x %v", errors.New("file not found"))),
			packType: []string{"*errors.errorString"},
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			pack := Unpack(v.err)

			pt := make([]string, len(pack))
			for i, p := range pack {
				pt[i] = fmt.Sprintf("%T", p.Cause())
			}

			if strings.Join(pt, ",") != strings.Join(v.packType, ",") {
				t.Fatalf("bad result for unpack, want: %v, but: %+v", strings.Join(v.packType, ","), strings.Join(pt, ","))
			}
		})
	}

	return
}
