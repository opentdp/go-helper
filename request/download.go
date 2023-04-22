package request

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func DownloadWithProgress(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// 获取文件大小
	fileSize := resp.ContentLength

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "download-*")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// 设置缓冲区大小，以减少内存使用
	bufferSize := 32 * 1024 // 32KB
	buf := make([]byte, bufferSize)

	// 追踪下载进度
	var downloaded int64
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Printf("\rDownloaded %d of %d bytes (%.2f%%)", downloaded, fileSize, float64(downloaded)/float64(fileSize)*100)
		}
	}()

	// 下载文件
	for {
		n, err := io.ReadFull(resp.Body, buf)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return "", err
		}
		downloaded += int64(n)

		_, err = tempFile.Write(buf[:n])
		if err != nil {
			return "", err
		}

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
	}

	ticker.Stop()
	fmt.Printf("\rDownloaded %d of %d bytes (100.00%%)\n", downloaded, fileSize)

	return tempFile.Name(), nil

}
