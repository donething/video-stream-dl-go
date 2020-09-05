package main

import (
	"bufio"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dohttp"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

// 视频保存的基路径
var basePath = "D:/Data/videos"
var client = dohttp.New(60*time.Second, false, false)

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
	addr := inputArgs()
	// 获取视频列表的文本
	text, err := getStreamText(addr)
	if err != nil {
		log.Printf("获取视频列表出错：%s ==>\n%s\n", err, text)
		return
	}
	// 解析文本内的视频地址
	names, format, err := ParsesHlsLinks(text)
	if err != nil {
		log.Fatalf("解析hls链接出错：%s ==>\n%s\n", err, text)
	}
	// 创建下载目录
	folder := strings.ReplaceAll(addr, "\\", "/")
	folder = addr[strings.LastIndex(folder, "/"):]
	if strings.Contains(folder, ".") {
		folder = folder[:strings.Index(folder, ".")]
	}
	// 初始化下载目录
	initDLDir(folder)

	// 待下载的总视频数量
	TotalCount = len(names)
	log.Printf("开始下载视频片段\n")
	for _, name := range names {
		// .m38u中的链接可能为绝对下载地址和相对地址，当为相对下载地址时，需要和addr合并
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
	UIWriter.Stop()

	// 处理下载结果
	// 有任务下载失败
	failed := TotalCount - DoneCount
	if failed != 0 {
		log.Fatalf("\n部分视频下载失败。总任务数：%d，失败的个数：%d\n", len(names), failed)
	}
	err = Combine(basePath, "", format)
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
func inputArgs() string {
	// 参数中包含了url
	if len(os.Args) == 2 {
		return os.Args[1]
	}

	reader := bufio.NewReaderSize(os.Stdin, 65536)
	log.Printf("请输入视频流的地址：")
	url, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("输入的视频流的地址有误：%s\n", err)
	}
	return strings.TrimSpace(url)
}

// 获取视频列表
func getStreamText(addr string) (string, error) {
	if strings.Index(addr, "http") == 0 {
		headers := map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"}
		text, err := client.GetText(addr, headers)
		if err != nil {
			return text, fmt.Errorf("根据地址获取视分段出错：%w", err)
		}
		return text, nil
	} else {
		bs, err := dofile.Read(addr)
		if err != nil {
			return "", fmt.Errorf("读取本地文件出错：%w", err)
		}
		return string(bs), nil
	}
}
