package main

import (
	"crypto/md5"
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for fileName := range jobs {

		data := []byte("содержимое файла " + fileName)
		hash := md5.Sum(data)

		results <- fmt.Sprintf("Файл %s -> MD5: %x", fileName, hash)
	}
}

func main() {
	files := []string{"doc1.txt", "doc2.txt", "doc3.txt", "doc4.txt", "doc5.txt"}

	jobs := make(chan string, len(files))
	results := make(chan string, len(files))

	var wg sync.WaitGroup

	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs) 

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println(res)
	}
}