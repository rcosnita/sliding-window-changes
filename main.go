package main

import (
	"fmt"
	"math/rand"
	"sliding-window-changes/sliding"
	"sync"
	"time"
)

func generateLoggers(slidingWindow sliding.SlidingWindow) {
	for {
		numOfChanges := rand.Int() % 1000

		fmt.Printf("[START] Generate %d num of changes in the sliding window", numOfChanges)
		fmt.Println()

		for idx := 0; idx < numOfChanges; idx++ {
			err := slidingWindow.Collect(sliding.SlidingWindowLogTaskChange{
				ComponentId:  "pon-c",
				LoggerFqdn:   "myLogger",
				CreationTime: time.Now().UnixMilli(),
				OldLogLevel:  "debug",
				NewLogLevel:  "info",
			})
			if err != nil {
				// TODO(rcosnita) make this smarter ...
				panic(err)
			}
		}

		fmt.Printf("[FINISH] Generate %d num of changes in the sliding window", numOfChanges)
		fmt.Println()
		time.Sleep(750 * time.Millisecond)
	}
}

func processAsync(slidingWindow sliding.SlidingWindow) {
	processor := sliding.NewSlidingWindowProcess(slidingWindow, "https://my-llm/")
	processor.Process()
}

func main() {
	slidingWindow := sliding.NewSlidingWindow()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		generateLoggers(slidingWindow)
		wg.Done()
	}()

	go func() {
		processAsync(slidingWindow)
		wg.Done()
	}()

	wg.Wait()

	fmt.Println("Experiment finished successfully ...")
}
