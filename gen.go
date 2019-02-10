package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

func main2() {
	day := time.Now()

	args := os.Args[1:]

	if len(args) < 1 {
		panic("Usage: \n\tgen 10\n\tgen 20 dirs")
	}

	count, _ := strconv.Atoi(args[0])
	mode := "files"
	if len(args) > 1 {
		mode = args[1]
	}

	fmt.Printf("Generate %d %s...\n", count, mode)
	os.Mkdir("test", 0755)

	for i := 0; i < count; i++ {
		// folder := fmt.Sprintf("")
		if mode == "dirs" {
			os.MkdirAll(path.Join("test", day.Format("2006-01-02")), 0755)
		} else {
			os.Create(path.Join("test", "backup_"+day.Format("2006-01-02")+".zip"))
		}
		day = day.Add(-time.Hour * 24)
	}
}
