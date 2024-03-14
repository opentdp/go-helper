package command

import (
	"errors"
	"os"
	"strings"
)

type ExecPayload struct {
	Name          string `note:"脚本名称"`
	CommandType   string `note:"脚本类型 BAT|POWERSHELL|SHELL"`
	Gb18030ToUtf8 bool   `note:"是否转换输出文本编码"`
	Username      string `note:"执行脚本的用户名"`
	WorkDirectory string `note:"脚本工作目录"`
	Content       string `note:"脚本内容"`
	Timeout       uint   `note:"超时时间"`
}

func Exec(data *ExecPayload) (string, error) {

	var (
		err error
		tmp string
		bin string
		arg []string
	)

	switch data.CommandType {
	case "BAT":
		tmp, err = newScript(data.Content, ".bat")
		arg = []string{"/c", "CALL", tmp}
		bin = "cmd.exe"
	case "POWERSHELL":
		tmp, err = newScript(data.Content, ".ps1")
		arg = []string{"-File", tmp}
		bin = "powershell.exe"
	case "SHELL":
		if strings.HasPrefix(data.Content, "#!/") {
			tmp, err = newScript(data.Content, "")
			arg = []string{}
			bin = tmp
		} else {
			arg = []string{"-c", data.Content}
			bin = "sh"
			tmp = "-"
		}
	case "EXEC":
		arg = strings.Fields(data.Content)
		bin, arg = arg[0], arg[1:]
		tmp = "-"
	default:
		err = errors.New("不支持此类脚本")
	}

	if err != nil || tmp == "" {
		return "", err
	}

	if tmp != "-" {
		defer os.Remove(tmp)
	}

	return execScript(bin, arg, data)

}
