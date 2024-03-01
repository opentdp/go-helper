package filer

import (
	"os"
)

// 判断是否目录
func IsDir(p string) bool {

	info, err := os.Stat(p)
	if err != nil {
		return false
	}

	return info.IsDir()

}

// 判断是否链接
func IsLink(p string) bool {

	info, err := os.Lstat(p)
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeSymlink != 0

}

// 判断文件是否存在
func Exists(p string) bool {

	if _, err := os.Stat(p); err != nil {
		return !os.IsNotExist(err)
	}

	return true

}
