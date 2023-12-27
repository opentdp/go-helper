//go:build windows

package filer

import (
	"fmt"
	"os"
)

func getFileOwner(fileInfo os.FileInfo) (string, string, error) {

	return "", "", fmt.Errorf("Not implemented")

}
