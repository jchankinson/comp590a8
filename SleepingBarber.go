package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	maxCustomers = 10
)

var (
	waitingRoom      = make(chan int, 5) 
	barberChair      = make(chan int, 1)
	wg               sync.WaitGroup
	noMoreCustomers bool 
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Barber Shop Simulation")

	wg.Add(1)
	go barber()
	time.Sleep(1 * time.Second) 

	for i := 1; i <= maxCustomers; i++ {
		wg.Add(1)
		go customer(i)
		if i == maxCustomers {
			noMoreCustomers = true
		}
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second) 
	}

	wg.Wait() 

	fmt.Println("Simulation complete")
}

func barber() {
	defer wg.Done()

	for {
		select {
		case customerID := <-barberChair:
			fmt.Printf("Barber is cutting hair for customer %d\n", customerID)
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second) 
			fmt.Printf("Barber finished cutting hair for customer %d\n", customerID)
		default:
			select {
			case customerID := <-waitingRoom:
				fmt.Printf("Barber is cutting hair for customer %d\n", customerID)
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
				fmt.Printf("Barber finished cutting hair for customer %d\n", customerID)
			default:
				if noMoreCustomers {
					fmt.Println("Barber is done for the day, no more customers")
					return
				}
				fmt.Println("Barber is sleeping")
				time.Sleep(500 * time.Millisecond) 
			}
		}
	}
}

func customer(id int) {
	defer wg.Done()

	fmt.Printf("Customer %d arrived at the barber shop\n", id)

	select {
	case barberChair <- id:
		fmt.Printf("Customer %d is waiting in the waiting room\n", id)
		return
	default:
		select {
		case waitingRoom <- id:
			fmt.Printf("Customer %d is waiting in the waiting room\n", id)
		default:
			fmt.Printf("Barber shop is full, Customer %d is leaving\n", id)
			return
		}
	}
}