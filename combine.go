package main

import (
	"bytes"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var files = "files.txt"

// 合并视频文件
// dst 目标目录
// filter 过滤视频，只通过普通的contains()判断
// format 视频格式，含点号，如".mp4"
func Combine(dst string, filter string, format string) error {
	log.Printf("开始列出文件，并保存到文件中\n")
	namesText := ""
	names := listFiles(dst, filter)
	sort.Strings(names)
	for _, name := range names {
		namesText += fmt.Sprintf("file '%s'\n", name)
	}

	_, err := dofile.Write([]byte(namesText), path.Join(dst, files), os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("保存文件名列表出错：%s\n", err)
		return err
	}

	log.Printf("使用FFmpeg合并视频文件\n")
	cmd := exec.Command("ffmpeg", "-f", "concat", "-i", "files.txt",
		"-c", "copy", "output"+format)
	cmd.Dir = dst
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Printf("执行合并文件命令时出错：%s == >%s\n", err, stderr.String())
		return err
	}
	return nil
}

// 获取视频目录的文件列表
func listFiles(dst string, filter string) []string {
	names := make([]string, 0, 20)
	err := filepath.Walk(dst, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Name() == files || !strings.Contains(info.Name(), filter) {
			return nil
		}
		names = append(names, info.Name())
		return nil
	})
	if err != nil {
		log.Printf("扫描视频目录(%s)出错：%s\n", dst, err)
		return []string{}
	}
	return names
}
