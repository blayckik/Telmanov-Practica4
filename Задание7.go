package main

import (
	"fmt"
	"sync"
)

type Request struct {
	Action string   
	Key    string    
	Value  string     
	Result chan string 
}

func stateManager(requests <-chan Request) {

	state := make(map[string]string)

	for req := range requests {
		switch req.Action {
		case "SET":
			state[req.Key] = req.Value
			req.Result <- "OK"
		case "GET":
			val, ok := state[req.Key]
			if !ok {
				req.Result <- "Ключ не найден"
			} else {
				req.Result <- val
			}
		}
	}
}

func main() {
	requests := make(chan Request)

	go stateManager(requests)

	var wg sync.WaitGroup

	keys := []string{"user", "role", "theme"}
	values := []string{"Алексей", "Admin", "Dark"}

	for i := 0; i < len(keys); i++ {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			resChan := make(chan string)
			requests <- Request{Action: "SET", Key: k, Value: v, Result: resChan}
			<-resChan 
		}(keys[i], values[i])
	}

	wg.Wait() 

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			resChan := make(chan string)
			requests <- Request{Action: "GET", Key: k, Result: resChan}
			val := <-resChan 
			fmt.Printf("Запрос ключа '%s' -> Значение: %s\n", k, val)
		}(key)
	}

	wg.Wait() 
}