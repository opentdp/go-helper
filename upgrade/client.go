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

func Downloader(rq *UpdateInfo) (io.ReadCloser, error) {

	resp, err := http.Get(rq.Package)

	if err != nil {
		return nil, fmt.Errorf("get package failed (%s)", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("get package failed (http status %d)", resp.StatusCode)
	}

	if strings.HasSuffix(rq.Package, ".gz") && resp.Header.Get("Content-Encoding") != "gzip" {
		return gzip.NewReader(resp.Body)
	}

	return resp.Body, nil

}
