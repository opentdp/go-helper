package upgrade

import (
	"bytes"
	"crypto"
	"errors"
	"os"
	"time"
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

	// Use to mark the temporary file.
	time string
}

func (o *Options) getPath() (string, error) {

	if o.TargetPath == "" {
		p, err := os.Executable()
		if err != nil {
			return "", err
		}
		o.TargetPath = p
	}

	return o.TargetPath, nil

}

func (o *Options) getMode() os.FileMode {

	if o.TargetMode == 0 {
		return 0755
	}

	return o.TargetMode

}

func (o *Options) getTimeString() string {

	if o.time == "" {
		o.time = time.Now().Format("20060102150405")
	}

	return o.time

}

func (o *Options) verifyChecksum(payload []byte) error {

	h := o.Hash
	if h == 0 {
		h = crypto.SHA256
	}

	if !h.Available() {
		return errors.New("requested hash function not available")
	}

	hash := h.New()
	hash.Write(payload) // guaranteed not to error
	checksum := hash.Sum([]byte{})

	if !bytes.Equal(o.Checksum, checksum) {
		return errors.New("updated file has wrong checksum")
	}

	return nil

}
