package main

import (
	"github.com/donething/utils-go/dohttp"
	"log"
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
)

// 需要初始化
func init() {
	// 启动goroutine来完成工作
	// 方便需要等待所有worker工作完毕才能推出
	WG.Add(workerCount)
	for gr := 1; gr <= workerCount; gr++ {
		// 可以go worker(TasksCh, gr)，而不能go func(){worker(TasksCh, gr)}
		go worker(TasksCh, gr)
	}
	log.Println("工作goroutine已准备就绪")
}

// 工作
func worker(tasks chan string, worker int) {
	defer WG.Done()
	for {
		task, ok := <-tasks
		if !ok {
			log.Printf("Worker:%d 已完成下载工作", worker)
			break
		}
		log.Printf("Worker:%d 收到任务，进行工作", worker)
		// 实际工作
		// 保存文件的路径
		name := task[strings.LastIndex(task, "/")+1:]
		if strings.Contains(name, "?") {
			name = name[:strings.Index(name, "?")]
		}
		p := path.Join(basePath, name)

		_, _, err := client.DownFile(task, p, false, nil)
		if err != nil && err != dohttp.ErrFileExist {
			log.Printf("下载视频片段失败：%v\n", err)
			continue
		}
		doneMutex.Lock()
		DoneCount++
		doneMutex.Unlock()
		log.Printf("Worker:%d 已下载完该任务[%d/%d]\n", worker, DoneCount, TotalCount)
	}
}
