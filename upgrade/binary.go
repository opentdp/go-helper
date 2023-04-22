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
	"path/filepath"
	"runtime"
	"strings"

	"github.com/open-tdp/go-helper/request"
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

// PrepareBinary reads the new binary content from io.Reader and performs the following actions:
//  1. If configured, applies the contents of the update io.Reader as a binary patch.
//  2. If configured, computes the checksum of the executable and verifies it matches.
//  3. If configured, verifies the signature with a public key.
//  4. Creates a new file, /path/to/.target.new with the TargetMode with the contents of the updated file

func PrepareBinary(update io.Reader, opts Options) error {

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

	// get the directory the executable exists in
	updateDir := filepath.Dir(targetPath)
	filename := filepath.Base(targetPath)

	// Copy the contents of newbinary to a new executable file
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, opts.getMode())
	if err != nil {
		return err
	}

	defer fp.Close()

	_, err = io.Copy(fp, bytes.NewReader(newBytes))
	if err != nil {
		return err
	}

	// if we don't call fp.Close(), windows won't let us move the new executable
	// because the file will still be "in use"
	fp.Close()

	return nil

}

// CommitBinary moves the new executable to the location of the current executable or opts.TargetPath
// if specified. It performs the following operations:
//  1. Renames /path/to/target to /path/to/.target.old
//  2. Renames /path/to/.target.new to /path/to/target
//  3. If the final rename is successful, deletes /path/to/.target.old, returns no error. On Windows,
//     the removal of /path/to/target.old always fails, so instead Apply hides the old file instead.
//  4. If the final rename fails, attempts to roll back by renaming /path/to/.target.old
//     back to /path/to/target.
//
// If the roll back operation fails, the file system is left in an inconsistent state where there is
// no new executable file and the old executable file could not be be moved to its original location.
// In this case you should notify the user of the bad news and ask them to recover manually. Applications
// can determine whether the rollback failed by calling RollbackError, see the documentation on that function
// for additional detail.

func CommitBinary(opts Options) error {

	// get the directory the file exists in
	targetPath, err := opts.getPath()
	if err != nil {
		return err
	}

	updateDir := filepath.Dir(targetPath)
	filename := filepath.Base(targetPath)
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))

	// this is where we'll move the executable to so that we can swap in the updated replacement
	oldPath := opts.OldSavePath
	removeOld := opts.OldSavePath == ""
	if removeOld {
		oldPath = filepath.Join(updateDir, fmt.Sprintf(".%s.old", filename))
	}

	// delete any existing old exec file - this is necessary on Windows for two reasons:
	// 1. after a successful update, Windows can't remove the .old file because the process is still running
	// 2. windows rename operations fail if the destination file already exists
	_ = os.Remove(oldPath)

	// move the existing executable to a new file in the same directory
	err = os.Rename(targetPath, oldPath)
	if err != nil {
		return err
	}

	// move the new exectuable in to become the new program
	err = os.Rename(newPath, targetPath)

	if err != nil {
		// move unsuccessful
		//
		// The filesystem is now in a bad state. We have successfully
		// moved the existing binary to a new location, but we couldn't move the new
		// binary to take its place. That means there is no file where the current executable binary
		// used to be!
		// Try to rollback by restoring the old binary to its original path.
		rerr := os.Rename(oldPath, targetPath)
		if rerr != nil {
			return &ErrRollback{err, rerr}
		}

		return err
	}

	// move successful, remove the old binary if needed
	if removeOld {
		os.Remove(oldPath)
	}

	return nil

}

// RollbackError takes an error value returned by Apply and returns the error, if any,
// that occurred when attempting to roll back from a failed update. Applications should
// always call this function on any non-nil errors returned by Apply.
//
// If no rollback was needed or if the rollback was successful, RollbackError returns nil,
// otherwise it returns the error encountered when trying to roll back.

func RollbackError(err error) error {

	if err == nil {
		return nil
	}

	if er, ok := err.(*ErrRollback); ok {
		return er.rollbackErr
	}

	return nil

}
