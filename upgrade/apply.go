package upgrade

import (
	"github.com/minio/selfupdate"
	"github.com/open-tdp/go-helper/logman"
)

func Apply(rq *RequesParam) error {

	logman.Info(
		"Checking update",
		"version", rq.Version,
		"update_url", rq.UpdateUrl,
	)

	update, err := Downloader(rq)
	if err != nil {
		logman.Error("Check update", "result", err)
		return err
	}

	defer update.Close()

	err = selfupdate.Apply(update, selfupdate.Options{})
	if err != nil {
		logman.Error("Apply update failed", "error", err)
		return err
	}

	return nil

}
