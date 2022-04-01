package main

import (
	"fmt"
	"snowflake"
	"sync"
	"time"
)

func main() {
	now := time.Now()
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			worker, _ := snowflake.NewWorker(int64(idx))
			for i := 0; i < 4000000; i++ {
				worker.Generate()
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(time.Now().Sub(now).String())
}