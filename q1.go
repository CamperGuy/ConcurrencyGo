package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func dentist(wait <-chan chan int, dent <-chan chan int, comms chan<- string) {
	admissionQueue := make(chan chan int, 5)

	go func() {
		for {
			select {
			case patient := <-wait:
				admissionQueue <- patient
				break
			default:
				wakeMeUpPatient := <-dent
				wakeMeUpPatient <- -200
			}
		}
	}()

	go func() {
		for patient := range admissionQueue {
			patient <- -101
			treatmentTime := rand.Intn(6-1) + 1
			comms <- "Dentist: Treating a patient for " + strconv.Itoa(treatmentTime) + " seconds"
			time.Sleep(time.Duration(treatmentTime) * time.Second)
			patient <- -100
		}
	}()
}

func patient(wait chan<- chan int, dent chan<- chan int, id int, comms chan<- string) {
	self := make(chan int)
	dent <- self // wake up the dentist if they were asleep
	<-self       // Accept confirmation that dentist is awake

	wait <- self // Hey, I'd like to be treated
	go func() {
		for sleepMessage := range self {
			fmt.Println("Patient " + strconv.Itoa(id) + " about to be treated")
			if sleepMessage == -101 {
				<-self
				fmt.Println("   Patient " + strconv.Itoa(id) + " done\n")
			}
		}
	}()
}

func main() {
	n, m := 5, 5
	dent := make(chan chan int)    // creates a synchronous channel
	wait := make(chan chan int, n) // creates an asynchronous
	comms := make(chan string, 50) // creates an asynchronous logging channel

	go func() {
		for message := range comms {
			fmt.Println(message)
		}
	}()

	// channel of size n
	go dentist(wait, dent, comms)
	time.Sleep(3 * time.Second)
	comms <- strconv.Itoa(m) + " patients"
	for i := 0; i < m; i++ {
		go patient(wait, dent, i, comms)
		time.Sleep(1 * time.Second)
	}
	comms <- "Main: Everything initialised"
	time.Sleep(20 * time.Second)
}
