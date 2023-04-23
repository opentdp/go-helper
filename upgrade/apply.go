package upgrade

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/open-tdp/go-helper/logman"
)

func Apply(rq *RequesParam) error {

	logger := logman.Named("upgrade")

	logger.Info(
		"checking update",
		"version", rq.Version,
		"url", rq.Server,
	)

	resp, err := CheckVersion(rq)
	if err != nil {
		logger.Error("check update failed", "error", err)
		return err
	}

	if !strings.HasPrefix(resp.Package, "https://") {
		logger.Info("no need to update", "resp", resp)
		return ErrNoUpdate
	}

	updater, err := Downloader(resp.Package)
	if err != nil {
		logger.Error("download binary failed", "error", err)
		return err
	}

	defer updater.Close()

	opts := Options{}

	if err = PrepareBinary(updater, opts); err != nil {
		logger.Error("prepare binary failed", "error", err)
		return err
	}

	if err = CommitBinary(opts); err != nil {
		logger.Error("apply update failed", "error", err)
		if RollbackError(err) != nil {
			logger.Error("failed to rollback from bad update")
		}
		return err
	}

	return nil

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
