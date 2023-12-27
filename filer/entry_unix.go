//go:build !windows

package filer

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
)

func getFileOwner(fileInfo os.FileInfo) (string, string, error) {

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return "", "", fmt.Errorf("Not a syscall.Stat_t")
	}

	usr, err := user.LookupId(fmt.Sprint(stat.Uid))
	if err != nil {
		return "", "", err
	}

	grp, err := user.LookupGroupId(fmt.Sprint(stat.Gid))
	if err != nil {
		return "", "", err
	}

	return usr.Username, grp.Name, nil

}
