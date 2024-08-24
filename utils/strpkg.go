package utils

import (
	"strings"
	"unicode"
)

func ExtractNumbers(input string) string {
	var sb strings.Builder
	for _, char := range input {
		if unicode.IsDigit(char) {
			sb.WriteRune(char)
		}
	}
	return sb.String()
}
