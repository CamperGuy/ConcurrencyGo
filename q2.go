package main

/*
Comment for part b
*/

import (
	"fmt"
	"time"
)

func dentist(hwait chan chan int, lwait <-chan chan int, dent <-chan chan int) {
	timer := time.NewTimer(m * time.Millisecond)
}

func patient(wait chan chan int, dent chan chan int, id int) {

}

func main() {
	dent := make(chan chan int)
	hwait := make(chan chan int, 100)
	lwait := make(chan chan int, 5)
	go dentist(hwait, lwait, dent)
	high := 10
	low := 3

	for i := low; i < high; i++ {
		go patient(hwait, dent, i)
	}

	for i := 0; i < low; i++ {
		go patient(lwait, dent, i)
	}

	time.Sleep(50 * time.Second)
	fmt.Println("Done")
}
