package main

/*
Identify a possibility of deadlock that affects part 2 but not part 3?
Justify

The dentist in Part 2 could deadlock when Line 50 is chosen.
	The timer would run out and looking for a low priority patient.
	If there is no low priority patient present, but only a high
	priority patient, the dentist would deadlock.
In Part 3 the Dentist would not deadlock as this would not be possible
	The assistant now also just keeps cycling when either waiting rooms are empty
*/

/*
	Math/rand, depending on your installation of Go might need to be lowercase
	I found this:
	macOS under Go 1.15.5: Uppercase
		Installed via Homebrew
	Windows under Go 1.15.5: Lowercase
		Installed via provided installer
	Please adjust for your system as required. I could not find a universal work-around
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

	go func() {
		for {
			select {
			// Set priority for the timer to be chosen when it has run out
			case <-lowTimer.C:
				select {
				case lowPatient := <-lwait:
					fmt.Println("Assistant: Timeout -> Low Priority Patient")
					wait <- lowPatient
					lowTimer.Reset(duration) // Restart the timer
				default: // If there is a timeout and the queue is empty just keep cycling
					lowTimer.Reset(duration)
					continue
				}
			default:
				select {
				case highPatient := <-hwait:
					fmt.Println("Assistant: High Priority Patient")
					wait <- highPatient
				default:
					select {
					case lowPatient := <-lwait:
						fmt.Println("Assistant: Low Priority Patient")
						wait <- lowPatient
					default: // If the queue is empty just keep cycling
						continue
					}
				}
			}
		}
	}()
}

func dentist(wait chan chan int, dent <-chan chan int) {
	for {
		waitingPatient := <-dent
		waitingPatient <- -200
		patient := <-wait
		patient <- -101
		treatmentTime := rand.Intn(6-1) + 1
		fmt.Println("Dentist: Treating patient for " + strconv.Itoa(treatmentTime) + " seconds\n")
		time.Sleep(time.Duration(treatmentTime) * time.Second)
		patient <- -100
	}
}

func patient(wait chan<- chan int, dent chan<- chan int, id int) {
	self := make(chan int)
	dent <- self
	<-self
	wait <- self
	fmt.Println("Patient " + strconv.Itoa(id) + ": waiting")
	go func() {
		for sleepMessage := range self {
			fmt.Println("Patient " + strconv.Itoa(id) + ": getting treatment")
			if sleepMessage == -101 {
				<-self
				fmt.Println("Patient " + strconv.Itoa(id) + ": done")
			}
		}
	}()
}

func main() {
	dent := make(chan chan int)
	hwait := make(chan chan int, 100)
	lwait := make(chan chan int, 100)
	wait := make(chan chan int)

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
