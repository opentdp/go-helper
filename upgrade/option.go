package upgrade

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Options struct {
	// TargetPath defines the path to the file to update.
	// The emptry string means 'the executable file of the running program'.
	TargetPath string

	// Create TargetPath replacement with this file mode. If zero, defaults to 0755.
	TargetMode os.FileMode

	// Checksum of the new binary to verify against. If nil, no checksum or signature verification is done.
	Checksum []byte

	// Use this hash function to generate the checksum. If not set, SHA256 is used.
	Hash crypto.Hash

	// Store the old executable file at this path after a successful update.
	// The empty string means the old executable file will be removed after the update.
	OldSavePath string
}

// CheckPermissions determines whether the process has the correct permissions to
// perform the requested update. If the update can proceed, it returns nil, otherwise
// it returns the error that would occur if an update were attempted.
func (o *Options) CheckPermissions() error {

	// get the directory the file exists in
	path, err := o.getPath()
	if err != nil {
		return err
	}

	fileDir := filepath.Dir(path)
	fileName := filepath.Base(path)

	// attempt to open a file in the file's directory
	newPath := filepath.Join(fileDir, fmt.Sprintf(".%s.check-perm", fileName))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, o.getMode())
	if err != nil {
		return err
	}
	fp.Close()

	_ = os.Remove(newPath)
	return nil

}

func (o *Options) getPath() (string, error) {

	if o.TargetPath == "" {
		return executable()
	} else {
		return o.TargetPath, nil
	}

}

func (o *Options) getMode() os.FileMode {

	if o.TargetMode == 0 {
		return 0755
	}
	return o.TargetMode

}

func (o *Options) getHash() crypto.Hash {

	if o.Hash == 0 {
		o.Hash = crypto.SHA256
	}
	return o.Hash

}

func (o *Options) verifyChecksum(updated []byte) error {

	checksum, err := checksumFor(o.getHash(), updated)
	if err != nil {
		return err
	}

	if !bytes.Equal(o.Checksum, checksum) {
		return errors.New("updated file has wrong checksum")
	}

	return nil

}

func checksumFor(h crypto.Hash, payload []byte) ([]byte, error) {

	if !h.Available() {
		return nil, errors.New("requested hash function not available")
	}

	hash := h.New()
	hash.Write(payload) // guaranteed not to error

	return hash.Sum([]byte{}), nil

}

func executable() (string, error) {

	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Abs(ex)

}
