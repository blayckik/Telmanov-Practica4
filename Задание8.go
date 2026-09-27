package main

import (
	"fmt"
	"sync"
	"time"
)
func serviceInstance(id int, requests <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for req := range requests {
		fmt.Printf("Экземпляр сервиса #%d обработал запрос: %s\n", id, req)
		time.Sleep(100 * time.Millisecond)
	}
}

func roundRobinBalancer(incoming <-chan string, serviceChans []chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	currentIndex := 0
	numServices := len(serviceChans)

	for req := range incoming {
		serviceChans[currentIndex] <- req
		currentIndex = (currentIndex + 1) % numServices
	}

	for _, ch := range serviceChans {
		close(ch)
	}
}

func main() {
	const numServices = 3 
	const totalRequests = 10 

	incomingRequests := make(chan string, totalRequests)

	serviceChans := make([]chan string, numServices)
	for i := 0; i < numServices; i++ {
		serviceChans[i] = make(chan string, 5)
	}

	var wgServices sync.WaitGroup
	var wgBalancer sync.WaitGroup

	for i := 0; i < numServices; i++ {
		wgServices.Add(1)
		go serviceInstance(i+1, serviceChans[i], &wgServices)
	}

	wgBalancer.Add(1)
	go roundRobinBalancer(incomingRequests, serviceChans, &wgBalancer)
	fmt.Println("Отправка входящих запросов в балансировщик:")
	for r := 1; r <= totalRequests; r++ {
		reqText := fmt.Sprintf("Запрос №%d", r)
		incomingRequests <- reqText
	}
	close(incomingRequests) 
	wgBalancer.Wait()
	wgServices.Wait()

	fmt.Println("Все запросы успешно распределены и обработаны!")
}