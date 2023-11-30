package filer

import (
	"io/fs"
	"os"
)

// 列出目录中的所有文件
func List(dir string) ([]fs.FileInfo, error) {

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var list []fs.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		list = append(list, info)
	}

	return list, nil

}

// 判断文件是否存在
func Exists(path string) bool {

	if _, err := os.Stat(path); err != nil {
		return !os.IsNotExist(err)
	}

	return true

}
