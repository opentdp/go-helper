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
	IsLink  bool        // 是否是链接
	IsDir   bool        // 是否是目录
	Data    []byte      // 文件数据
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
			IsLink:  info.Mode()&os.ModeSymlink != 0,
			IsDir:   IsDir(filepath.Join(dir, file.Name())),
		})
	}

	return list, nil

}

// 获取文件信息和内容
func Detail(path string, read bool) (*FileInfo, error) {

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	detail := &FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode().Perm(),
		ModTime: info.ModTime().Unix(),
		IsLink:  IsLink(path),
		IsDir:   info.IsDir(),
	}

	if read && !info.IsDir() {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		detail.Data = data
	}

	return detail, nil

}

// 写入文件内容，目录不存在时自动创建
func Write(path string, data []byte) error {

	// 创建目录
	if dir := filepath.Dir(path); !Exists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	// 写入内容
	return os.WriteFile(path, data, 0644)

}
