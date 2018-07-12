package main

import (
	"bufio"
	"os"
	"io"
	"fmt"
	"strings"
	"strconv"
	"net/http"
)

func main() {
	fmt.Println("========= Downloader Started! =========")

	files := os.Args[1:]
	if len(files) == 0 {
	} else {
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ImageDownloader: %v\n", err)
				continue
			}

			scanLines(f)
			f.Close()
		}
	}
}

func scanLines(f *os.File) {
	input := bufio.NewScanner(f)
	ch := make(chan string)

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

			go download(url, size, code, ch)
			fmt.Println(<-ch)
		}
	}
}

func download(url string, size int64, _ string, ch chan<- string) {
	// 检查文件是否已经存在
	splited := strings.Split(url, "/")
	filename := "downloads/" + splited[len(splited) - 1]
	f0, err := os.Open(filename)
	if err == nil {
		stat, err := f0.Stat()
		if err == nil && stat.Size() == size {
			ch <- fmt.Sprintf("%s already downloaded!", filename)
			return
		}
	}

	resp, err := http.Get(url)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	nbytes, err := io.Copy(f, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	if size != nbytes {
		ch <- fmt.Sprintf("%s size incorrect: %d -> %d", url, size, nbytes)
		os.Remove(filename)
		return
	}

	ch <- fmt.Sprintf("%s download to %s successful!", url, filename)
}
