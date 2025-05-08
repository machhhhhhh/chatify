package utils

import (
	"regexp"
)

func IsCorrectFormatEmail(email string) bool {
	var rule string = `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	var regex *regexp.Regexp = regexp.MustCompile(rule)
	return regex.MatchString(email) == true
}

func IsCorrectFormatPassword(password string) bool {
	var IsHasUpper bool = regexp.MustCompile(`[A-Z]`).MatchString(password)
	var IsHasLower bool = regexp.MustCompile(`[a-z]`).MatchString(password)
	var IsHasDigit bool = regexp.MustCompile(`[0-9]`).MatchString(password)
	return IsHasUpper == true && IsHasLower == true && IsHasDigit == true
}
