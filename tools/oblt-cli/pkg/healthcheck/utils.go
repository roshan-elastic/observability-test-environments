package healthcheck

import (
	"regexp"
)

func IsNotBlank(s string) bool {
	if s == "" {
		return false
	}
	if regexp.MustCompile(`^\s+$`).MatchString(s) {
		return false
	}
	return true
}
