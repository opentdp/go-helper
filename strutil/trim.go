package strutil

import (
	"strings"
)

// 删除文本中每一行的公共前导空白
func Dedent(text string) string {

	lines := strings.Split(text, "\n")
	minIndent := -1

	for _, line := range lines {
		trimLine := strings.TrimLeft(line, " \t")
		if trimLine == "" {
			continue // 跳过空行
		}

		indent := len(line) - len(trimLine)
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent > 0 {
		for i, line := range lines {
			if line = strings.TrimLeft(line, " \t"); len(line) > 0 {
				lines[i] = line
			}
		}
	}

	// 移除首尾的空行
	start, end := 0, len(lines)
	for start < end && lines[start] == "" {
		start++
	}
	for end > start && lines[end-1] == "" {
		end--
	}

	return strings.Join(lines[start:end], "\n")

}
