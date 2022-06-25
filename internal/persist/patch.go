package persist

import (
	"strings"
)

func JsonPointerToJsonPath(s string) string {
	var b strings.Builder
	b.WriteRune('$')

	parts := strings.Split(s, "/")
	for _, part := range parts {
		if len(part) > 0 {
			b.WriteRune('.')
			b.WriteString(part)
		}
	}

	return b.String()
}
