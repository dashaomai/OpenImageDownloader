package main

import (
	"bufio"
	"flag"
	"os"
	"io"
	"fmt"
	"strings"
	"strconv"
	"net/http"
)

type DownloadTask struct {
	url string
	size int64
	code string
}

func main() {
	fmt.Println("========= Downloader Started! =========")

	files := os.Args[1:]
	if len(files) == 0 {
	} else {
		var n = flag.Int("n", 32, "number of threads")
		flag.Parse()

		var threadNumber = *n

		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				// 文件不存在，可能是其它控制参数
				// fmt.Fprintf(os.Stderr, "ImageDownloader: %v\n", err)
				continue
			}

			scanLines(f, threadNumber)
			f.Close()
		}
	}
}

func scanLines(f *os.File, threadNumber int) {
	input := bufio.NewScanner(f)

	fmt.Printf("%d thread(s) should created.\n", threadNumber)

	// 创建指定数量的 channel 用于和一一对应的协程进行通讯
	ch := make(chan DownloadTask)
	for i := 0; i<threadNumber; i++ {
		go download(ch)
	}

	for input.Scan() {
		line := input.Text()
		args := strings.Split(line, "\t")
		if len(args) == 3 {
			url := args[0]
			size, err := strconv.ParseInt(args[1], 10, 64)
			code := args[2]

			if err != nil {
				fmt.Println(err)
				continue
			}

			ch <- DownloadTask{url, size, code}
		}
	}
}

// 分配任务的进程
func download(ch <-chan DownloadTask) {
	for task := range ch {
		// fmt.Printf("Process: %v\n", task)
		download2(task)
		// fmt.Printf("Processed: %v\n", task)
	}
}

func download2(task DownloadTask) {
	url := task.url
	size := task.size

	// 检查文件是否已经存在
	splited := strings.Split(url, "/")
	filename := "downloads/" + splited[len(splited) - 1]
	f0, err := os.Open(filename)
	if err == nil {
		stat, err := f0.Stat()
		if err == nil && stat.Size() == size {
			fmt.Printf("%s already downloaded!\n", filename)
			return
		}
	}

	resp, err := http.Get(url)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	nbytes := resp.ContentLength
	if size != nbytes {
		fmt.Printf("%s size incorrect: %d -> %d\n", url, size, nbytes)
		return
	}

	f, err := os.Create(filename)
	if (f != nil) {
		defer f.Close()
	}

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	_, err2 := io.Copy(f, resp.Body)
	if err2 != nil {
		fmt.Printf("while reading %s: %v\n", url, err2)
		return
	}

	fmt.Printf("%s download successful!\n", url)
}
