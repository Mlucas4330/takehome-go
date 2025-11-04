package validator

import (
	"regexp"
	"strings"
)

func ValidateRG(rg string) bool {
	rg = strings.ReplaceAll(rg, ".", "")
	rg = strings.ReplaceAll(rg, "-", "")
	rg = strings.ToUpper(rg)

	if len(rg) < 5 || len(rg) > 20 {
		return false
	}

	matched, _ := regexp.MatchString(`^[A-Z0-9]+$`, rg)
	return matched
}
