package command

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/opentdp/go-helper/logman"
	"github.com/opentdp/go-helper/strutil"
)

func newScript(code string, ext string) (string, error) {

	tf, err := os.CreateTemp("", "tmp-*"+ext)

	if err != nil {
		return "", errors.New("创建临时文件失败")
	}

	defer tf.Close()

	// 替换换行符
	lineEnding := detectLineEnding(code)
	if runtime.GOOS == "windows" && lineEnding == "unix" {
		code = strings.ReplaceAll(code, "\n", "\r\n")
	} else if runtime.GOOS != "windows" && lineEnding == "windows" {
		code = strings.ReplaceAll(code, "\r\n", "\n")
	} else if runtime.GOOS != "windows" && lineEnding == "mac" {
		code = strings.ReplaceAll(code, "\r", "\n")
	}

	// 写入临时文件
	if _, err = tf.WriteString(code); err != nil {
		return "", errors.New("写入临时文件失败")
	}

	// 赋予执行权限
	if runtime.GOOS != "windows" {
		tf.Chmod(0755)
	}

	return tf.Name(), nil

}

func execScript(bin string, arg []string, data *ExecPayload) (string, error) {

	logman.Debug("执行应用程序", "bin", bin, "arg", arg)

	// 超时时间
	timeout := time.Duration(data.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	// 运行脚本文件
	cmd := exec.CommandContext(ctx, bin, arg...)

	// 设置执行目录
	if data.WorkDirectory != "" {
		cmd.Dir = data.WorkDirectory
	}

	// 获取输出信息
	ret, err := cmd.CombinedOutput()
	str := string(ret)

	// 转换文本编码
	if data.Gb18030ToUtf8 {
		str = strutil.Gb18030ToUtf8(str)
	}

	return str, err

}

func detectLineEnding(code string) string {

	if strings.Contains(code, "\r\n") {
		return "windows"
	}
	if strings.Contains(code, "\n") {
		return "unix"
	}
	if strings.Contains(code, "\r") {
		return "mac"
	}

	return "unknown"

}
