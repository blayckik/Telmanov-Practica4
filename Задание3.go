package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const totalRequests = 15

	tick := time.Tick(200 * time.Millisecond)

	var wg sync.WaitGroup

	fmt.Println("Начало обработки запросов к API:")
	startTime := time.Now()

	for i := 1; i <= totalRequests; i++ {
		wg.Add(1)

		<-tick

		go func(reqID int) {
			defer wg.Done()

			fmt.Printf("[%s] Запрос #%d успешно обработан\n",
				time.Since(startTime).Truncate(time.Millisecond), reqID)
		}(i)
	}

	wg.Wait()
	fmt.Println("Все запросы обработаны!")
}