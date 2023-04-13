package upgrade

type RequesParam struct {
	UpdateUrl string `note:"检测地址"`
	Version   string `note:"当前版本"`
}

type UpdateInfo struct {
	BinaryUrl string `note:"下载地址"`
	Message   string `note:"错误信息"`
	Version   string `note:"最新版本"`
}
