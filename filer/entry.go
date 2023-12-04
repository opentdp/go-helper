package filer

import (
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name    string      // 文件名
	Size    int64       // 字节大小
	Mode    os.FileMode // 权限，如 0777
	ModTime int64       // 修改时间，Unix时间戳
	Symlink string      // 链接的真实路径，软链接时有效
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
		fp := filepath.Join(dir, file.Name())
		list = append(list, &FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			Mode:    info.Mode().Perm(),
			ModTime: info.ModTime().Unix(),
			Symlink: Readlink(fp),
			IsDir:   IsDir(fp),
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
		Symlink: Readlink(path),
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

	if dir := filepath.Dir(path); !Exists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return os.WriteFile(path, data, 0644)

}

// 获取软链接的真实路径
func Readlink(path string) string {

	if IsLink(path) {
		if rp, err := os.Readlink(path); err == nil {
			return rp
		}
	}
	return ""

}
