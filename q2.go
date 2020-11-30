package main

/*
Comment for part b
*/

import (
	"fmt"
	"time"
	"Math/rand"
	"strconv"
)

/*
--- Task ---
Dentist will check that there are no high-prioirty patients before starting on a low priorty one
Need to account for starvation: use aging technique
Move a patient from lwait to hwait whenever m milliseconds have passed since the last read
 from lwait

*/
func dentist(hwait chan chan int, lwait <-chan chan int, dent <-chan chan int) {
	highQueue := make(chan chan int, 10)
	lowQueue := make(chan chan int, 10)

	// Add people into their corresponding queue
	go func() {
		for {
			timer := time.NewTimer(200 * time.Millisecond)

			select {
			case highPatient := <-hwait:
				highQueue <- highPatient
				
			case <-timer.C:
				select {
				case lowPatient := <-lwait:
					lowQueue <- lowPatient
				default:
					wakeMeUpPatient := <- dent // Dentist to fall asleep
					wakeMeUpPatient <- -200
				}
			}
		}
	}()

	// Treat them!
	go func() {
		for {
			patientChan:= make (chan chan int)

			select {
			case patientChan = <- highQueue:
				patientChan <- -100
			default patientChan = <- lowQueue:
				patientChan <- - 100
			} 	

			treatmentTime := rand.Intn(6-1) + 1
			time.Sleep(time.Duration(treatmentTime) * time.Second)
			patientChan <- -101
		}
	}()
}

// Patient must also work for Q1
func patient(wait chan<- chan int, dent chan<- chan int, id int) {
	self := make(chan int)
	dent <- self // wake up the dentist if they were asleep
	<- dent // Accept confirmation that dentist is awake

	wait <- self // Hey, I'd like to be treated
	go func(){
		for sleepMessage := range self{
			if sleepMessage == -101{
				wakeUpMessage := <- self
				fmt.PrintlnI("Patient " + strconv.Itoa(id) + " done")
			}
		}
	}

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
