package upgrade

import (
	"strings"

	"github.com/minio/selfupdate"
	"github.com/open-tdp/go-helper/logman"
)

func Apply(rq *RequesParam) error {

	logger := logman.Named("updater")

	logger.Info(
		"checking update",
		"version", rq.Version,
		"url", rq.UpdateUrl,
	)

	info, err := CheckVersion(rq)
	if err != nil {
		logger.Error("check update failed", "error", err)
		return err
	}

	if !strings.HasPrefix(info.BinaryUrl, "https://") {
		logger.Info("no need to update", "info", info)
		return ErrNoUpdate
	}

	updater, err := Downloader(info)
	if err != nil {
		logger.Error("prepare updater failed", "error", err)
		return err
	}

	defer updater.Close()

	err = selfupdate.Apply(updater, selfupdate.Options{})
	if err != nil {
		logger.Error("apply update failed", "error", err)
		return err
	}

	return nil

}
