package request

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"

	"github.com/cheggaaa/pb/v3"
)

func Download(url string, showProgress bool) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 默认读取器为响应体
	reader := resp.Body

	// 如果需要显示下载进度
	if showProgress {
		bar := pb.StartNew(int(resp.ContentLength))
		bar.Set(pb.Bytes, true) //自动换为合适的字节单位
		reader = bar.NewProxyReader(reader)
		defer bar.Finish()
	}

	// 检查是否使用gzip压缩
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return "", err
		}
		defer reader.Close()
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "tdp-*")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// 将下载的数据复制到临时文件
	_, err = io.Copy(tempFile, reader)
	if err != nil {
		return "", err
	}

	// 返回临时文件的名称
	return tempFile.Name(), nil

}
