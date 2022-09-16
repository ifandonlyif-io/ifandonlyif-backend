package util

import (
	"strings"
)

func SnakeCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	//snake_case to camelCase

	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	camelCase = replaceFirstRune(camelCase, strings.ToLower(string(camelCase[0])))
	return
}

func replaceFirstRune(str, replacement string) string {
	return string([]rune(str)[:0]) + replacement + string([]rune(str)[1:])
}
