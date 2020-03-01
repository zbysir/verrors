package verrors

import (
	"fmt"
	"sort"
	"strings"
)

// formatPackErrors简单的打印错误, 方便临时查看
func formatPackErrors(ps PackErrors) string {
	var s strings.Builder
	for _, v := range ps {
		if s.Len() != 0 {
			s.WriteString("\n")
		}

		s.WriteString("- " + formatPackError(v))
	}

	return s.String()
}

// formatPackError简单的打印错误, 方便临时查看
func formatPackError(e PackError) string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%v", e.err))
	values := e.GetAll()

	if len(values) != 0 {
		var keys []string
		for k := range values {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		kv := make([]string, len(keys))
		for i, k := range keys {
			kv[i] = fmt.Sprintf("%s = %v", k, values[k])
		}

		s.WriteString(fmt.Sprintf(" [ %s ]", strings.Join(kv, "; ")))
	}

	return s.String()
}
