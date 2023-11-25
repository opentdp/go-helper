package upgrade

import (
	"errors"
)

type RequesParam struct {
	Server  string `note:"更新服务器"`
	Version string `note:"当前版本"`
}

type UpdateInfo struct {
	Type    string `note:"更新方式"`
	Error   string `note:"错误信息"`
	Message string `note:"提示信息"`
	Release string `note:"更新说明"`
	Version string `note:"最新版本"`
	Package string `note:"下载地址"`
}

type ErrRollback struct {
	error          // original error
	Rollback error // error encountered while rolling back
}

var ErrNoUpdate = errors.New("no update")
