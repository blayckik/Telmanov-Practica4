package main

import (
	"fmt"
	"math/rand"
	"time"
)

func searchInSource(sourceName string, query string, resultChan chan<- string) {
	delay := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(delay)

	result := fmt.Sprintf("Результат по запросу '%s' найден в источнике [%s] за %v", query, sourceName, delay)
	resultChan <- result
}

func main() {
	rand.Seed(time.Now().UnixNano())

	query := "Go concurrency"
	sources := []string{"База данных №1", "База данных №2", "Внешний API", "Кэш-сервер"}

	resultChan := make(chan string, len(sources))

	for _, source := range sources {
		go searchInSource(source, query, resultChan)
	}

	firstResult := <-resultChan

	fmt.Println("Первый успешный результат:")
	fmt.Println(firstResult)
}