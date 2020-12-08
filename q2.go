package main

/*
Assume Go does not have fair semantics
Where could Starvation occur in Part 1?
Justify answer.
*/

import (
	"Math/rand"
	"fmt"
	"strconv"
	"time"
)

func dentist(hwait chan chan int, lwait <-chan chan int, dent <-chan chan int) {
	highQueue := make(chan chan int, 10)
	lowQueue := make(chan chan int, 10)

	// Add people into their corresponding queue
	go func() {
		for {
			select {
			case highPatient := <-hwait:
				highQueue <- highPatient
			case lowPatient := <-lwait:
				lowQueue <- lowPatient
			default:
				wakeMeUpPatient := <-dent // Dentist to fall asleep
				wakeMeUpPatient <- -200   // Send awake confirmation
			}
		}
	}()

	// Treatment Routine
	go func() {
		// Every 8 seconds, treat a low patient to prevent starvation
		duration := 8 * time.Second
		lowTimer := time.NewTimer(duration)

		patientChan := make(chan int)
		for {
			select {
			// Set priority for the timer to be chosen when it has run out
			case <-lowTimer.C:
				fmt.Println("--- Dentist choosing Low due to timeout ---")
				patientChan = <-lowQueue
				patientChan <- -101
				lowTimer.Reset(duration) // Restart the timer
			default:
				// Then normally, prefer highQueue patients over lowQueue patients
				select {
				case patientChan = <-highQueue:
					fmt.Println("--- Emergency ---")
					patientChan <- -101
				default:
					fmt.Println("--- Regular ---")
					patientChan = <-lowQueue
					patientChan <- -101
				}
			}

			// Treat the patient (go to sleep) and send them a wake up signal
			treatmentTime := rand.Intn(6-1) + 1
			time.Sleep(time.Duration(treatmentTime) * time.Second)
			patientChan <- -100
		}
	}()
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
