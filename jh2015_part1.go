package main

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

func dentist(wait <-chan chan int, dent <-chan chan int) {
	admissionQueue := make(chan chan int, 5)

	go func() {
		for {
			// Add patients to the queue
			select {
			case patient := <-wait:
				admissionQueue <- patient
			default:
				/* Sleep until we receive a wakeup and confirm that we received it
				As per Brief */
				wakeMeUpPatient := <-dent
				wakeMeUpPatient <- -200
			}
		}
	}()

	go func() {
		// Every patient coming in, send them a sleep signal, wait, wake them up
		for patient := range admissionQueue {
			patient <- -101
			treatmentTime := rand.Intn(5) + 1
			fmt.Println("Dentist is busy for " + strconv.Itoa(treatmentTime) + " seconds")
			time.Sleep(time.Duration(treatmentTime) * time.Second)
			patient <- -100
		}
	}()
}

func patient(wait chan<- chan int, dent chan<- chan int, id int) {
	self := make(chan int)
	// Wake up dentist and receive confirmation message
	dent <- self
	<-self

	// Add myself to queue
	wait <- self
	fmt.Println("Patient " + strconv.Itoa(id) + ": waiting")
	go func() {
		// Stay in a waiting state until we receive a treatment start signal
		for sleepMessage := range self {
			fmt.Println("  Patient " + strconv.Itoa(id) + ": getting treatment")
			if sleepMessage == -101 {
				<-self // And wait until we are woken up
				fmt.Println("    Patient " + strconv.Itoa(id) + ": done\n")
			}
		}
	}()
}

func main() {
	n, m := 5, 5
	dent := make(chan chan int)    // creates a synchronous channel
	wait := make(chan chan int, n) // creates an asynchronous

	// channel of size n
	go dentist(wait, dent)
	time.Sleep(3 * time.Second)
	fmt.Println(strconv.Itoa(m) + " patients")
	for i := 0; i < m; i++ {
		go patient(wait, dent, i)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(30 * time.Second)
}
