package main

import (
	"errors"
	"github.com/donething/utils-go/dohttp"
	"strings"
	"time"
)

var client = dohttp.New(60*time.Second, false, false)

// 从.m38u文本中返回 文件列表和视频格式
func ParsesHlsLinks(text string) ([]string, string, error) {
	names := make([]string, 0)
	format := ""
	lines := strings.Split(text, "\n")
	// 不是.m38u格式的文件
	if len(lines) < 0 || lines[0] != "#EXTM3U" {
		return nil, "", errors.New("无法解析.m38u文件")
	}

	for _, line := range lines {
		if strings.TrimSpace(line) != "" && line[0] != '#' {
			names = append(names, line)
		}
	}
	// 获取视频格式，FFmpeg合并需要指定视频后缀
	if len(names) > 0 && strings.Contains(names[0], ".") {
		format = names[0][strings.LastIndex(names[0], "."):]
		if strings.Contains(format, "?") {
			format = format[:strings.Index(format, "?")]
		}
	}
	return names, format, nil
}
