package upgrade

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"

	"github.com/open-tdp/go-helper/request"
)

func CheckVersion(rq *RequesParam) (*UpdateInfo, error) {

	info := &UpdateInfo{}

	body, err := request.TextGet(rq.UpdateUrl, request.H{
		"app-version":      rq.Version,
		"app-runtime-os":   runtime.GOOS,
		"app-runtime-arch": runtime.GOARCH,
	})

	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(body), &info)

	if err != nil {
		return info, err
	}
	if info.Message != "" {
		return info, errors.New(info.Message)
	}
	if info.BinaryUrl == "" {
		return info, errors.New("get package url failed")
	}

	return info, nil

}

func Downloader(rq *UpdateInfo) (io.ReadCloser, error) {

	resp, err := http.Get(rq.BinaryUrl)

	if err != nil {
		return nil, fmt.Errorf("get package failed (%s)", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("get package failed (http status %d)", resp.StatusCode)
	}

	if strings.HasSuffix(rq.BinaryUrl, ".gz") && resp.Header.Get("Content-Encoding") != "gzip" {
		return gzip.NewReader(resp.Body)
	}

	return resp.Body, nil

}
