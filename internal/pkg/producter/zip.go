package main

import (
	"collector/internal/pkg/producter/collectzip"
	"collector/pkg/cache"
	"collector/pkg/config"
	"collector/pkg/logger"
	"fmt"
	"time"
)

func main() {
	z := collectzip.New(logger.New("./", logger.LOGALL), cache.New("./cache", "db"), config.New("./collecter.json"), nil)
	z.Run()

	for {
		select {
		case file := <-z.QueueFiles:
			z.DoneOneFiles(file)
			fmt.Println("QueueFiles:", file)

			fmt.Println()
			time.Sleep(2 * time.Second)
		}

	}
}
