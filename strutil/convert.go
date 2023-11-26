package strutil

import (
	"strconv"
	"strings"
)

// 字符串转为 int

func ToInt(str string) int {

	v, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return v

}

// 字符串转为 uint

func ToUint(str string) uint {

	v, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return uint(v)

}

// 字符串首字母大写

func FirstUpper(s string) string {

	if s == "" {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]

}

// 字符串首字母小写

func FirstLower(s string) string {

	if s == "" {
		return s
	}

	return strings.ToLower(s[:1]) + s[1:]

}
