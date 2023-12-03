package command

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/opentdp/go-helper/strutil"
)

func newScript(code string, ext string) (string, error) {

	tf, err := os.CreateTemp("", "go-*."+ext)

	if err != nil {
		return "", errors.New("创建临时文件失败")
	}

	defer tf.Close()

	if _, err = tf.WriteString(code); err != nil {
		return "", errors.New("写入临时文件失败")
	}

	if runtime.GOOS != "windows" {
		tf.Chmod(0755)
	}

	return tf.Name(), nil

}

func execScript(bin string, arg []string, data *ExecPayload) (string, error) {

	timeout := time.Duration(data.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	cmd := exec.CommandContext(ctx, bin, arg...)

	if data.WorkDirectory != "" {
		cmd.Dir = data.WorkDirectory
	}

	ret, err := cmd.CombinedOutput()
	str := string(ret)

	if runtime.GOOS == "windows" {
		str = strutil.Gb18030ToUtf8(str)
	}

	return str, err

}
