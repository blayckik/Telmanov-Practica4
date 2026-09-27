package main

import (
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		square := job * job
		results <- square
	}
}

func main() {
	const numJobs = 10
	const numWorkers = 3


	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("Результаты работы воркеров (возведение в квадрат):")
	for res := range results {
		fmt.Println(res)
	}
}