package bvutils

import "strings"

// SeparateCamel 按照自定符号分隔驼峰
func SeparateCamel(name, separator string) string {
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteString(separator)
		}

		result.WriteRune(r | ' ')
	}
	return result.String()
}

// UcFirst 字符串首字母大写
func UcFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// LcFirst 字符串首字母小写
func LcFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
