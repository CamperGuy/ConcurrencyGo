package main

/*
Notes:
Treatment should be done as a handshake

Waiting room can fit every patient
- Don't worry about behaviour where this is an issue
Use message passing
Asynchronous channel (wait) for the waiting queue
Synchronous channel (dent) to model the state (sleeping/awake) of the DENTIST
	! Dentist should fall asleep when reading on dent, not while reading on wait !
Use synchronous channels (one per patient) to model the state (sleeping/awake) of that patient
	a) patient (function) creates this as a fresh channel (chan int)
	b) This channel will be added to the wait (chan chan int) if the patient needs to stay in
	   the waiting room
	c) Patient sleeps while waiting to read on their fresh channel when:
		1) waiting in the queue
		2) having some treatment done
Treatment should take a random time
*/

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

/* 	--- Task ---
Check for people in the waiting room
If there is a patient, call the first one in
During treatment, dentist is active and the patient is sleeping (blocked)
When finshed with treatment, patient is woken up
Dentist checks for next person in the waiting room
*/
func dentist(wait <-chan chan int, dent <-chan chan int, comms chan<- string) {
	busy := false
	admissionQueue := make(chan chan int, 5)

	// Direct Admission
	go func() {
		for {
			patientChan := <-dent
			if busy {
				patientChan <- -203 // Send rejection signal. They will now request the waiting room
				comms <- "  Dentist: Direct admission rejected Patient"
			} else {
				patientChan <- -202 // Send acceptance signal
				comms <- "  Dentist: Direct admission accepted Patient"
			}

			admissionQueue <- patientChan
		}
	}()

	// Waiting Room
	go func() {
		for {
			patientChan := <-wait

			patientChan <- -201 // Send rejection signal to send them to sleep
			comms <- "  Dentist: Patient admitted to waiting room"

			admissionQueue <- patientChan // Add them to our queue
		}
	}()

	// Treatment
	go func() {
		for {
			patientChan := <-admissionQueue
			busy = true
			patientChan <- -100 // Send them to sleep
			treatmentTime := rand.Intn(6-1) + 1
			comms <- "  Dentist: Treating a patient for " + strconv.Itoa(treatmentTime) + " seconds"
			time.Sleep(time.Duration(treatmentTime) * time.Second)
			patientChan <- -101
			busy = false
		}
	}()
}

/* 	--- Task ---``
Check if the dentist is busy. If so, fall asleep (block?)
If dentist is sleeping:
	Wake up dentist and self fall asleep
	- We will be woken up when the treatment is completed
If dentist is busy with another patient:
	Go to the waiting room
	Sleep
	When woken up, go to sleep after treatment started
*/
func patient(wait chan<- chan int, dent chan<- chan int, id int, comms chan<- string) {
	self := make(chan int)
	dent <- self
	go func() {
		for {
			response := <-self
			// comms <- "Patient " + strconv.Itoa(id) + ": received response " + strconv.Itoa(response)

			if response == -203 {
				// Direct admission unsuccessful
				comms <- "Patient " + strconv.Itoa(id) + ": going to waiting room"
				wait <- self
			} else if response == -201 {
				// Go to sleep to be woken up in the waiting room
				comms <- "Patient " + strconv.Itoa(id) + ": falling asleep in the waiting room"
			} else if response == -100 {
				comms <- "Patient " + strconv.Itoa(id) + ": I will receive treatment"
				if <-self == -101 {
					comms <- "Patient " + strconv.Itoa(id) + ": I'm done!\n"
				}
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
	comms <- "\n"
	for i := 0; i < m; i++ {
		go patient(wait, dent, i, comms)
		time.Sleep(1 * time.Second)
		comms <- "\n"
	}
	comms <- "Main: Everything initialised"
	time.Sleep(1 * time.Second)
}
