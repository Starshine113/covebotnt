package etc

import (
	"strings"
)

// TrimPrefixesSpace trims all given prefixes as well as whitespace from the given string
func TrimPrefixesSpace(s string, prefixes ...string) string {
	for _, prefix := range prefixes {
		s = strings.TrimPrefix(s, prefix)
	}
	s = strings.TrimSpace(s)
	return s
}

// HasAnyPrefix checks if the string has *any* of the given prefixes
func HasAnyPrefix(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
