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

func Downloader(rq *RequesParam) (io.ReadCloser, error) {

	info, err := CheckVersion(rq)

	if err != nil {
		return nil, fmt.Errorf("check version (%s)", err)
	}

	resp, err := http.Get(info.BinaryUrl)

	if err != nil {
		return nil, fmt.Errorf("get request failed (%s)", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("get request failed (status code %d)", resp.StatusCode)
	}

	if strings.HasSuffix(info.BinaryUrl, ".gz") && resp.Header.Get("Content-Encoding") != "gzip" {
		return gzip.NewReader(resp.Body)
	}

	return resp.Body, nil

}

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
		return info, errors.New("get update url failed")
	}

	return info, nil

}
