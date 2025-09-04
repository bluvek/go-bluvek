package bvutils

import (
	"fmt"
	"reflect"
	"strings"
)

// MakePasswd 生成密码
func MakePasswd(pwd, salt string) string {
	return Md5Encode(pwd + salt)
}

// GetMapValue 获取map的值
func GetMapValue[T any](m map[string]interface{}, key string) T {
	var zero T

	value, exists := m[key]
	if exists {
		v := reflect.ValueOf(value)
		if v.Type().ConvertibleTo(reflect.TypeOf(zero)) {
			return v.Convert(reflect.TypeOf(zero)).Interface().(T)
		}
	}

	return zero
}

type MapSupportedTypes interface {
	string | int64 | float64 | bool
}

// GetMapSpecificValue 获取map的特定类型的值, 相较于GetMapValue, 不用每次反射获取值
func GetMapSpecificValue[T MapSupportedTypes](m map[string]interface{}, key string) T {
	var zero T

	value, exists := m[key]
	if exists {
		if v, ok := value.(T); ok {
			return v
		}
	}

	var result any
	switch v := value.(type) {
	case float64:
		if _, ok := any(zero).(int64); ok {
			result = int64(v)
		} else if _, ok := any(zero).(bool); ok {
			result = v != 0
		} else {
			return zero
		}
	case int64:
		if _, ok := any(zero).(float64); ok {
			result = float64(v)
		} else if _, ok := any(zero).(bool); ok {
			result = v != 0
		} else {
			return zero
		}
	case string:
		if _, ok := any(zero).(bool); ok {
			lowerVal := strings.ToLower(v)
			if lowerVal == "true" || lowerVal == "1" {
				result = true
			} else if lowerVal == "false" || lowerVal == "0" {
				result = false
			} else {
				return zero
			}
		} else {
			result = v
		}
	case bool:
		if _, ok := any(zero).(string); ok {
			result = fmt.Sprintf("%v", v) // 转成 "true" / "false"
		} else {
			result = v
		}
	default:
		return zero
	}

	if finalValue, ok := result.(T); ok {
		return finalValue
	}

	return zero
}

// InArray 判断某个值是否在数组中
func InArray[T comparable](val T, slice []T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// 三元运算符
func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
