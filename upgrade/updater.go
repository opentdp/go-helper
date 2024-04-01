package upgrade

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"os"
)

type Updater struct {
	// 新版本。默认为 'new'
	NewVersion string
	// 旧版本。默认为 'old'
	OldVersion string
	// 要更新的文件的路径。默认为 '正在运行的文件'
	TargetPath string
	// 可执行文件的权限掩码。默认为 0755
	TargetMode os.FileMode
	// 要应用的新文件的路径，必须为可执行文件
	NewBinary string
	// 新二进制文件的SHA256校验和。默认不进行校验
	Checksum []byte
}

func (u *Updater) Init() error {

	if u.NewVersion == "" {
		u.NewVersion = "new"
	}

	if u.OldVersion == "" {
		u.OldVersion = "old"
	}

	if u.TargetPath == "" {
		p, err := os.Executable()
		if err != nil {
			return err
		}
		u.TargetPath = p
	}

	if u.TargetMode == 0 {
		u.TargetMode = 0755
	}

	u.NewBinary = u.TargetPath + "-" + u.NewVersion

	return nil

}

func (u *Updater) VerifyChecksum() error {

	if u.Checksum == nil {
		return nil
	}

	// 打开文件
	file, err := os.Open(u.NewBinary)
	if err != nil {
		return err
	}
	defer file.Close()

	// 计算校验和
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	// 检查校验和
	if !bytes.Equal(u.Checksum, hash.Sum(nil)) {
		return errors.New("updated file has wrong checksum")
	}

	return nil

}

func (u *Updater) CommitBinary() error {

	// check the checksum if needed
	if err := u.VerifyChecksum(); err != nil {
		return err
	}

	// set the newBinary permission
	if err := os.Chmod(u.NewBinary, u.TargetMode); err != nil {
		return err
	}

	// backup the old binary
	originFile := u.TargetPath + "-" + u.OldVersion
	if err := os.Rename(u.TargetPath, originFile); err != nil {
		return err
	}

	// move the new exectuable in to become the new program
	if err := os.Rename(u.NewBinary, u.TargetPath); err != nil {
		// Try to rollback by restoring the old binary to its original path.
		if er2 := os.Rename(originFile, u.TargetPath); er2 != nil {
			return ErrRollback{err, er2}
		}
		return err
	}

	// try to remove the old binary if needed
	os.Remove(originFile)

	return nil

}
