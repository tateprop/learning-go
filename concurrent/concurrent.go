package main

import (
	"fmt"
	"time"
)

func hello() {
	time.Sleep(2 * time.Second)
	fmt.Println("Hello world goroutine")
}

func main() {
	go hello()
	time.Sleep(1 * time.Second)
	fmt.Println("main function")
	time.Sleep(2 * time.Second)
}
