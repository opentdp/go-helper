package upgrade

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/opentdp/go-helper/request"
)

func CheckVersion(rq *RequesParam) (*UpdateInfo, error) {

	info := &UpdateInfo{}

	url := rq.Server
	url += "?ver=" + rq.Version
	url += "&os=" + runtime.GOOS
	url += "&arch=" + runtime.GOARCH

	body, err := request.Get(url, request.H{})
	if err != nil {
		return info, err
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		return info, err
	}

	if info.Error != "" {
		return info, errors.New(info.Error)
	}
	if info.Package == "" {
		return info, errors.New("get package url failed")
	}

	return info, nil

}

func Downloader(url string) (io.ReadCloser, error) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("get package failed (%s)", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("get package failed (http status %d)", resp.StatusCode)
	}

	if strings.HasSuffix(url, ".gz") && resp.Header.Get("Content-Encoding") != "gzip" {
		return gzip.NewReader(resp.Body)
	}

	return resp.Body, nil

}

// reads the new binary content from io.Reader and performs the following actions:
//  If configured, applies the contents of the update io.Reader as a binary patch.
//  If configured, computes the checksum of the executable and verifies it matches.
//  Creates a new file with the TargetMode with the contents of the updated file

func PrepareBinary(update io.Reader, opts *Options) error {

	// get target path
	targetPath, err := opts.getPath()
	if err != nil {
		return err
	}

	var newBytes []byte

	// no patch to apply, go on through
	if newBytes, err = io.ReadAll(update); err != nil {
		return err
	}

	// verify checksum if requested
	if opts.Checksum != nil {
		if err = opts.verifyChecksum(newBytes); err != nil {
			return err
		}
	}

	// Copy the contents of newbinary to a new executable file
	newPath := targetPath + "-new" + opts.getTimeString()
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, opts.getMode())
	if err != nil {
		return err
	}

	defer fp.Close()

	_, err = io.Copy(fp, bytes.NewReader(newBytes))
	return err

}

// moves the new executable to the location of the current executable or opts.TargetPath

func CommitBinary(opts *Options) error {

	// get the directory the file exists in
	targetPath, err := opts.getPath()
	if err != nil {
		return err
	}

	newPath := targetPath + "-new" + opts.getTimeString()
	oldPath := targetPath + "-old" + opts.getTimeString()

	// move the existing executable to a new file in the same directory
	if err = os.Rename(targetPath, oldPath); err != nil {
		return err
	}

	// move the new exectuable in to become the new program
	if err = os.Rename(newPath, targetPath); err != nil {
		// Try to rollback by restoring the old binary to its original path.
		if er2 := os.Rename(oldPath, targetPath); er2 != nil {
			return &ErrRollback{err, er2}
		}
		return err
	}

	// try to remove the old binary if needed
	os.Remove(oldPath)

	return nil

}

// takes an error value returned by Apply and returns the error

func RollbackError(err error) error {

	if err == nil {
		return nil
	}

	if er, ok := err.(*ErrRollback); ok {
		return er.rollbackErr
	}

	return nil

}
