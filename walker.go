package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"
)

func checkExecutable(mode fs.FileMode) bool {
	if mode&0111 == 0111 {
		return true
	}
	return false
}

func scanFiles(path string) {
	skip := false
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("scanFiles() => os.Open(%s) => %s", path, err)
	}
	files, err := file.ReadDir(-1)
	if err != nil {
		fmt.Printf("scanFiles() => file.ReadDir(-1) => %s", err)
	}
	for _, file := range files {
		filepath := fmt.Sprintf("%s/%s", path, file.Name())
		stat, err := os.Stat(filepath)
		if err != nil {
			fmt.Printf("scanFiles() => os.Stat(%s) => %s", filepath, err)
		}
		if stat.IsDir() || checkExecutable(stat.Mode()) {
			continue
		}
		lastTime, ok := filesMap[filepath]
		if ok {
			if lastTime != stat.ModTime() {
				fmt.Println("File changed:", filepath)
				filesMap[filepath] = stat.ModTime()
				if !skip {
					skip = true
					go executeCommands()
				}

			}
		} else {
			filesMap[filepath] = stat.ModTime()
		}
	}
	time.Sleep(time.Second / 10)
}
