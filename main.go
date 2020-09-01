package main

import (
	"errors"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"log"
	"os"
	"path"
	"strings"
)

// 视频保存的基路径
var basePath = `D:\Data\videos`

// 创建工作目录
func initDLDir(folder string) {
	basePath = path.Join(basePath, folder)
	err := os.MkdirAll(basePath, 0644)
	if err != nil {
		log.Fatalf("创建基目录(%s)出错：%s\n", basePath, err)
	}
	log.Printf("视频保存的路径为 %s\n", basePath)
}

func main() {
	// 视频流的地址
	addr, filter := inputArgs()

	// 获取视频列表的文本
	text, err := getStreamText(addr)
	if err != nil {
		log.Printf("获取视频列表出错：%s\n", err)
		return
	}
	names, format, err := ParsesHlsLinks(text)
	if err != nil {
		log.Fatalf("解析hls链接出错：%s ==> %s\n", err, text)
	}
	// 创建下载目录
	folder := strings.ReplaceAll(addr, "\\", "/")
	folder = addr[strings.LastIndex(folder, "/"):]
	if strings.Contains(folder, ".") {
		folder = folder[:strings.Index(folder, ".")]
	}
	initDLDir(folder)

	TotalCount = len(names)
	for _, name := range names {
		u := name
		if !strings.Contains(name, "//") {
			u = addr[:strings.LastIndex(addr, "/")] + "/" + name
		}
		// 传递数据
		TasksCh <- u
	}
	// 完成了下载工作，关闭通道
	close(TasksCh)
	WG.Wait()

	// 有任务下载失败
	failed := TotalCount - DoneCount
	if failed != 0 {
		log.Fatalf("有视频下载失败。总任务数：%d，失败的个数：%d\n", len(names), failed)
	}
	err = Combine(basePath, filter, format)
	if err != nil {
		log.Fatalf("合并视频出错：%s\n", err)
	}

	log.Printf("已完成任务")
	err = dofile.ShowInExplorer(basePath)
	if err != nil {
		log.Printf("显示文件夹出错：%s\n", err)
	}
}

// 输入、检查参数
// 返回输入的视频流的地址和过滤文件的字符串
func inputArgs() (url string, filter string) {
	if len(os.Args) == 3 {
		return os.Args[1], os.Args[2]
	}

	log.Printf("请输入视频流的地址，和过滤文件的字符（此可空），以空格分隔")
	_, err := fmt.Scanln(&url, &filter)
	if err != nil && !strings.Contains(err.Error(), "unexpected newline") {
		log.Fatalf("输入错误：%s\n", err)
		return
	}
	url = strings.TrimSpace(url)
	if url == "" {
		log.Fatalf("输入的视频流的地址参数有误：%s\n", url)
	}
	return
}

// 获取视频列表
func getStreamText(addr string) (string, error) {
	if strings.Index(addr, "http") == 0 {
		text, _, err := client.GetText(addr, nil)
		if err != nil {
			return "", errors.New("根据地址（%s）获取视分段出错：" + err.Error())
		}
		return text, nil
	} else {
		bs, err := dofile.Read(addr)
		if err != nil {
			return "", errors.New("读取本地文件出错：" + err.Error())
		}
		return string(bs), nil
	}
}
