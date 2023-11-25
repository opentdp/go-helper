package request

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"

	"github.com/cheggaaa/pb/v3"
)

func Download(url, target string, showProgress, isGzip bool) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 默认读取器
	reader := resp.Body

	// 显示下载进度
	if showProgress {
		bar := pb.StartNew(int(resp.ContentLength))
		bar.Set(pb.Bytes, true) //自动换为合适的字节单位
		reader = bar.NewProxyReader(reader)
		defer bar.Finish()
	}

	// 自动解压缩文件
	if isGzip || resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return "", err
		}
	}

	// 尝试关闭读取器
	defer reader.Close()

	// 返回文件的名称
	return SaveStream(reader, target)

}

func SaveStream(reader io.Reader, target string) (string, error) {

	var err error
	var writer *os.File

	// 创建目标文件
	if target != "" {
		writer, err = os.Create(target)
	} else {
		writer, err = os.CreateTemp("", "tdp-*")
	}
	if err != nil {
		return "", err
	}
	defer writer.Close()

	// 写入文件数据
	_, err = io.Copy(writer, reader)
	if err != nil {
		return "", err
	}

	return writer.Name(), nil

}
