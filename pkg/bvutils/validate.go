package bvutils

import "strconv"

// ValidatePasswd 校验密码是否一致
func ValidatePasswd(pwd, salt, passwd string) bool {
	return Md5Encode(pwd+salt) == passwd
}

type Number interface {
	~int | ~int32 | ~int64 | ~float64 | ~float32 | ~string
}

// IsValidNumber 判断是否是有效数字
func IsValidNumber[T Number](value T) bool {
	switch v := any(value).(type) {
	case int:
		return v > 0
	case int32:
		return v > 0
	case int64:
		return v > 0
	case float64:
		return v > 0
	case float32:
		return v > 0
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num > 0
		}
	}
	return false
}
