package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
)

func main() {
	fmt.Println("========= List Started! =========")

	files := os.Args[1:]
	if len(files) == 0 {
		fmt.Println("list.exe file1.tsv file2.tsv file3.tsv ...")
		os.Exit(0)
	}

	out_path := "download.list"
	out, err := os.Create(out_path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file %s: %v\n", out_path, err)
		os.Exit(0)
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "List: %v\n", err)
			continue
		}

		scanLines(f, out)

		f.Close()
	}
	out.Close()
}

func scanLines(f *os.File, out *os.File) {
	input := bufio.NewScanner(f)
	output := bufio.NewWriter(out)
	defer output.Flush()

	for input.Scan() {
		line := input.Text()
		args := strings.Split(line, "\t")
		if len(args) == 3 {
			url := args[0]

			_, err := output.WriteString(url + "\n")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Write: %v\nScan Broken\n", err)
				break
			}
		}
	}
}
