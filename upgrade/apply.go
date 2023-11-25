package upgrade

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/opentdp/go-helper/logman"
	"github.com/opentdp/go-helper/request"
)

func Apply(rq *RequesParam) error {

	logger := logman.Named("upgrade")

	logger.Info(
		"checking update",
		"version", rq.Version,
		"url", rq.Server,
	)

	// check update

	resp, err := CheckVersion(rq)
	if err != nil {
		logger.Error("check update failed", "error", err)
		return err
	}

	if !strings.HasPrefix(resp.Package, "https://") {
		logger.Info("no need to update", "resp", resp)
		return ErrNoUpdate
	}

	// init updater

	updater := &Updater{}
	updater.Init()

	// download binary

	_, err = request.Download(resp.Package, updater.NewBinary, true)
	if err != nil {
		logger.Error("download binary failed", "error", err)
		return err
	}

	// apply binary update

	if err = updater.CommitBinary(); err != nil {
		logger.Error("apply binary failed", "error", err)
		if _, ok := err.(*ErrRollback); ok {
			logger.Error("failed to rollback from bad update")
		}
		return err
	}

	return nil

}

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

func Restart() error {

	self, err := os.Executable()
	if err != nil {
		return err
	}

	args, env := os.Args, os.Environ()

	// Windows does not support exec syscall
	if runtime.GOOS == "windows" {
		cmd := exec.Command(self, args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env
		if err = cmd.Start(); err != nil {
			return err
		}
		os.Exit(0)
	}

	// Other OS
	return syscall.Exec(self, args, env)

}
