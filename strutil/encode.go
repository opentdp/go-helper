package strutil

import (
	"golang.org/x/text/encoding/simplifiedchinese"
)

// 转换编码 gb18030 -> utf-8

func Gb18030ToUtf8(s string) string {

	ret, err := simplifiedchinese.GB18030.NewDecoder().String(s)
	if err == nil {
		return string(ret)
	}

	return s

}
