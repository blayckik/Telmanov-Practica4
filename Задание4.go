package main

import (
	"fmt"
	"net/http"
	"sync"
)

func checkURLWorker(id int, urls <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("[Воркер %d] Ошибка при запросе к %s: %v\n", id, url, err)
			continue
		}
		fmt.Printf("[Воркер %d] URL: %s | Статус: %s\n", id, url, resp.Status)
		resp.Body.Close()
	}
}

func main() {
	urlList := []string{
		"https://pkg.go.dev",
		"https://go.dev",
		"https://archlinux.org",
		"https://postman-echo.com/status/200",
		"https://postman-echo.com/status/404",
		"https://httpbin.org/get",
	}

	const numWorkers = 3
	urls := make(chan string, len(urlList))

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go checkURLWorker(w, urls, &wg)
	}

	for _, url := range urlList {
		urls <- url
	}
	close(urls)
	wg.Wait()
	fmt.Println("Проверка всех URL завершена!")
}