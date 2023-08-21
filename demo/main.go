package main

import (
	"fmt"
	"time"
)

var globalValue int

func main() {
	startTime := time.Now()

	for i := 0; i < 10000; i++ {
		go action(i)
	}

	delta := time.Now().Sub(startTime)
	fmt.Printf("Result is %d, done in %.3fs.\n", globalValue, delta.Seconds())
}

func action(i int) {
	fmt.Print(i, " ")

	time.Sleep(time.Second)
}
