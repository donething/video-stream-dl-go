package main

import (
	"errors"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"github.com/gosuri/uilive"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

var (
	// 工作缓冲数量
	workerCount = 10

	// 创建缓冲通道来传递工作数据
	TasksCh    = make(chan string, workerCount)
	TotalCount = 0
	DoneCount  = 0
	doneMutex  sync.Mutex

	// 等待所有worker工作完毕
	WG sync.WaitGroup

	// 同一行显示输出
	UIWriter = uilive.New()
)

// 需要初始化
func init() {
	// 启动goroutine来完成工作
	// 方便需要等待所有worker工作完毕才能推出
	WG.Add(workerCount)
	for gr := 1; gr <= workerCount; gr++ {
		// 可以go worker(TasksCh, gr)，而不能go func(){worker(TasksCh, gr)}
		go worker(TasksCh)
	}
	log.Println("工作goroutine已准备就绪")

	UIWriter.Start()
}

// 工作
func worker(tasks chan string) {
	defer WG.Done()
	doneCount := 0
	for {
		task, ok := <-tasks
		if !ok {
			//log.Printf("Worker:%d 已完成下载工作", worker)
			break
		}
		//log.Printf("Worker:%d 收到任务，进行工作", worker)
		// 实际工作
		// 保存文件的路径
		name := task[strings.LastIndex(task, "/")+1:]
		if strings.Contains(name, "?") {
			name = name[:strings.Index(name, "?")]
		}
		dst := path.Join(basePath, name)

		_, err := client.Download(task, dst, false, nil)
		if err != nil && !errors.Is(err, dohttp.ErrFileExists) {
			log.Printf("下载视频片段(%s)失败：%v\n", task, err)
			_ = os.Remove(dst)
			continue
		}
		doneMutex.Lock()
		DoneCount++
		doneCount = DoneCount
		doneMutex.Unlock()
		_, _ = fmt.Fprintf(UIWriter, "共 %d 个视频，已下载 %d 个\n", TotalCount, doneCount)
	}
}
