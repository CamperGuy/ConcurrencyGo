package main

/*
Identify a possibility of deadlock that affects part 2 but not part 3?
Justify
*/

import (
	"Math/rand"
	"fmt"
	"strconv"
	"time"
)

func assistant(hwait chan chan int, lwait <-chan chan int, wait chan<- chan int, dent chan chan int) {
	duration := 8 * time.Second
	lowTimer := time.NewTimer(duration)
	for {
		select {
		case <-lowTimer.C:
			patient := <-lwait
			wait <- patient
			lowTimer.Reset(duration)
		default:
			select {
			case patient := <-hwait:
				wait <- patient
			default:
				patient := <-lwait
				wait <- patient
			}
		}
	}
}

func dentist(wait chan chan int, dent <-chan chan int) {
	for {
		waitingPatient := <-dent
		waitingPatient <- -200
		patient := <-wait
		patient <- -101
		treatmentTime := rand.Intn(6-1) + 1
		time.Sleep(time.Duration(treatmentTime) * time.Second)
		patient <- -100
	}
}

// See Part 1
func patient(wait chan<- chan int, dent chan<- chan int, id int) {
	self := make(chan int)
	dent <- self
	<-self
	wait <- self
	fmt.Println("Patient " + strconv.Itoa(id) + ": waiting")
	go func() {
		for sleepMessage := range self {
			fmt.Println("  Patient " + strconv.Itoa(id) + ": getting treatment")
			if sleepMessage == -101 {
				<-self
				fmt.Println("    Patient " + strconv.Itoa(id) + ": done\n")
			}
		}
	}()
}

func main() {
	dent := make(chan chan int)
	hwait := make(chan chan int, 100)
	lwait := make(chan chan int, 100)
	wait := make(chan chan int, 200)

	go dentist(wait, dent)
	go assistant(hwait, lwait, wait, dent)

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
