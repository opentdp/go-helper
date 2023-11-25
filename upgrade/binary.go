package upgrade

import (
	"bytes"
	"crypto"
	"errors"
	"io"
	"os"
	"time"
)

type Updater struct {
	// 要更新的文件的路径。默认为 '正在运行的文件'
	TargetPath string
	// 可执行文件的权限掩码。默认为 0755
	TargetMode os.FileMode
	// 要应用的新二进制文件的路径。此参数不能为空
	NewBinary string
	// 新二进制文件的SHA256校验和。默认不进行校验
	Checksum []byte
}

func (u *Updater) getTarget() (string, error) {

	if u.TargetPath == "" {
		p, err := os.Executable()
		if err != nil {
			return "", err
		}
		u.TargetPath = p
	}

	return u.TargetPath, nil

}

func (u *Updater) getTargetMode() os.FileMode {

	if u.TargetMode == 0 {
		u.TargetMode = 0755
	}

	return u.TargetMode

}

func (u *Updater) verifyChecksum() error {

	// 打开文件
	file, err := os.Open(u.NewBinary)
	if err != nil {
		return err
	}
	defer file.Close()

	// 计算校验和
	hash := crypto.SHA256.New()
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
	if err := u.verifyChecksum(); err != nil {
		return err
	}

	// get the directory the file exists in
	targetPath, err := u.getTarget()
	if err != nil {
		return err
	}

	// set the newBinary permission
	mode := u.getTargetMode()
	if err = os.Chmod(u.NewBinary, mode); err != nil {
		return err
	}

	// move the existing executable to a new file in the same directory
	originFile := targetPath + "-" + time.Now().Format("20060102150405")
	if err = os.Rename(targetPath, originFile); err != nil {
		return err
	}

	// move the new exectuable in to become the new program
	if err = os.Rename(u.NewBinary, targetPath); err != nil {
		// Try to rollback by restoring the old binary to its original path.
		if er2 := os.Rename(originFile, targetPath); er2 != nil {
			return ErrRollback{err, er2}
		}
		return err
	}

	// try to remove the old binary if needed
	os.Remove(originFile)

	return nil

}
