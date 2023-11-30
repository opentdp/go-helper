package filer

import (
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name    string      // 文件名
	Size    int64       // 字节大小
	Mode    os.FileMode // 权限，如：0777
	ModTime int64       // 修改时间，Unix时间戳
	IsDir   bool        // 是否是目录
	Content string      // 文本内容，仅用于文本文件
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
			Size:    info.Size(),
			Mode:    info.Mode().Perm(),
			ModTime: info.ModTime().Unix(),
			IsDir:   info.IsDir(),
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
		Size:    info.Size(),
		Mode:    info.Mode().Perm(),
		ModTime: info.ModTime().Unix(),
		IsDir:   info.IsDir(),
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

// 读取文本内容
func Read(path string) (string, error) {

	data, err := os.ReadFile(path)

	return string(data), err

}

// 写入文本内容
func Write(path string, content string) error {

	// 创建目录
	if dir := filepath.Dir(path); !Exists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	// 写入内容
	return os.WriteFile(path, []byte(content), 0644)

}

// 判断文件是否存在
func Exists(path string) bool {

	if _, err := os.Stat(path); err != nil {
		return !os.IsNotExist(err)
	}

	return true

}
