package filer

import (
	"os"
)

// 判断是否目录
func IsDir(path string) bool {

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()

}

// 判断是否链接
func IsLink(path string) bool {

	info, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeSymlink != 0

}

// 判断文件是否存在
func Exists(path string) bool {

	if _, err := os.Stat(path); err != nil {
		return !os.IsNotExist(err)
	}

	return true

}
