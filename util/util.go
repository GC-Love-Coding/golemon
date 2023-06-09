package util

import "unicode"

func Max(a, b int) int {
	if b < a {
		return a
	}

	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func IsSpace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t' || r == '\v' || r == '\r'
}

func IsAlphaNum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func AllMatch(str string, predicate func(r rune) bool) bool {
	for _, r := range str {
		if !predicate(r) {
			return false

		}
	}

	return true
}

// Check the input string is all upper case or not.
// Each rune inside the string will be checked by `unicode.IsUpper` method.
// Note that call this method on empty string will get true.
// And the time complexity is O(N) where N is approximately length of input string.
func IsUpper(str string) bool {
	return AllMatch(str, func(r rune) bool {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}

		return true
	})
}

func IsLower(str string) bool {
	return AllMatch(str, func(r rune) bool {
		if unicode.IsLetter(r) && !unicode.IsLower(r) {
			return false
		}

		return true
	})
}

func IsLetter(str string) bool {
	return AllMatch(str, unicode.IsLetter)
}

// Checks if a given string is a string literal enclosed in single quotes ('') or double quotes ("").
func IsStringLiteral(str string) bool {
	if len(str) < 2 {
		return false
	}

	return (str[0] == '"' && str[len(str)-1] == '"') || (str[0] == '\'' && str[len(str)-1] == '\'')
}
