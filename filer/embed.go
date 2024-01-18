package filer

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

// 释放 embed.FS
func ReleaseEmbedFS(efs embed.FS, root string) (string, error) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "wrest-*")
	if err != nil {
		return "", err
	}
	// 递归复制目录内容
	err = fs.WalkDir(efs, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 计算目标路径
		relPath, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(tempDir, relPath)
		// 如果是目录，则创建目录
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}
		// 如果是文件，则复制文件内容
		data, err := fs.ReadFile(efs, p)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, data, d.Type().Perm())
	})
	// 出错时清理临时目录
	if err != nil {
		os.RemoveAll(tempDir)
		return "", err
	}
	// 返回临时目录的路径
	return tempDir, nil
}
