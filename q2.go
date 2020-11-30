package main

/*
Comment for part b
*/

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Test
/*
--- Task ---
Dentist will check that there are no high-prioirty patients before starting on a low priorty one
Need to account for starvation: use aging technique
Move a patient from lwait to hwait whenever m milliseconds have passed since the last read
 from lwait

*/
func dentist(hwait chan chan int, lwait <-chan chan int, dent <-chan chan int, comms chan<- string) {
	highQueue := make(chan chan int, 20)
	lowQueue := make(chan chan int, 20)

	// Add people into their corresponding queue
	go func() {
		for {
			comms <- "Dentist asleep until woken up"
			<-dent
			timer := time.NewTimer(200 * time.Millisecond)
			select {

			case highPatient := <-hwait:
				comms <- "Dentist: taking on a HIGH priority patient"
				highQueue <- highPatient
			case <-timer.C:

				select {
				case lowPatient := <-lwait:
					comms <- "Dentist: taking on a LOW priority patient"
					lowQueue <- lowPatient
				default:
				}
			}
		}

		/*
			for {

				select {
				case highPatient := <-hwait: // Priority needed to check for timeouts
					comms <- "Dentist: added high patient to high queue"
					highQueue <- highPatient

				case <-timer.C: // If time has passed, check for low (so they could be treated instantly)
					comms <- "Dentist : Timer has passed"
					select {
					case lowPatient := <-lwait:
						lowQueue <- lowPatient
						comms <- "Dentist: added low patient to low queue"
					default: // If there is nobdoy waiting fall asleep
						comms <- "Dentist: Falling asleep and waiting to be woken up"
						wakeMeUpPatient := <-dent
						wakeMeUpPatient <- -200
					}
				}
				comms <- "Dentist: at the end of a cycle"
			}
		*/
	}()

	// 2 Queues, one with high that should be FIFO and done first
	// One with low who should be FIFO and done when high is done

	// Treatment
	go func() {
		for {
			select {
			case highPatient := <-highQueue:
				// Treat high Patient
				highPatient <- -100
				treatmentTime := rand.Intn(6-1) + 1
				comms <- "  Dentist: Treating a HIGH patient for " + strconv.Itoa(treatmentTime) + " seconds"
				time.Sleep(time.Duration(treatmentTime) * time.Second)
				highPatient <- -101
			case lowPatient := <-lowQueue:
				// Treat low Patient
				lowPatient <- -100
				treatmentTime := rand.Intn(6-1) + 1
				comms <- "  Dentist: Treating a LOW patient for " + strconv.Itoa(treatmentTime) + " seconds"
				time.Sleep(time.Duration(treatmentTime) * time.Second)
				lowPatient <- -101
			}
		}
	}()
}

// Patient must also work for Q1
func patient(wait chan<- chan int, dent chan<- chan int, id int, comms chan<- string) {
	self := make(chan int)
	dent <- self // wake up the dentist if they were asleep
	wait <- self // add myself to relevant queue
	comms <- "Patient " + strconv.Itoa(id) + ": Waiting for treament signal"

	for treatSignal := range self { // -100 has been received
		comms <- "Patient " + strconv.Itoa(id) + " received " + strconv.Itoa(treatSignal)
		<-self // Wait for second signal
		comms <- "Patient " + strconv.Itoa(id) + ": Done"
	} //test

}

func main() {
	comms := make(chan string, 100)

	go func() {
		for message := range comms {
			fmt.Println(message)
		}
	}()

	dent := make(chan chan int)
	hwait := make(chan chan int, 100)
	lwait := make(chan chan int, 5)
	go dentist(hwait, lwait, dent, comms)
	high := 10
	low := 3

	for i := low; i < high; i++ {
		go patient(hwait, dent, i, comms)
	}

	for i := 0; i < low; i++ {
		go patient(lwait, dent, i, comms)
	}

	time.Sleep(50 * time.Second)
	fmt.Println("Done")
}
