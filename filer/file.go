package filer

import (
	"os"
)

type FileInfo struct {
	Name    string
	IsDir   bool
	Size    int64
	Mode    os.FileMode
	ModTime int64
	Content string
}

// 列出目录中的所有文件
func List(dir string) ([]*FileInfo, error) {

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var list []*FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		list = append(list, &FileInfo{
			Name:    info.Name(),
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			Mode:    info.Mode().Perm(),
			ModTime: info.ModTime().Unix(),
		})
	}

	return list, nil

}

// 获取文件信息和文本内容
func Detail(path string, text bool) (*FileInfo, error) {

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	detail := &FileInfo{
		Name:    info.Name(),
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		Mode:    info.Mode().Perm(),
		ModTime: info.ModTime().Unix(),
	}

	if text && !info.IsDir() {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		detail.Content = string(content)
	}

	return detail, nil

}

// 判断文件是否存在
func Exists(path string) bool {

	if _, err := os.Stat(path); err != nil {
		return !os.IsNotExist(err)
	}

	return true

}
