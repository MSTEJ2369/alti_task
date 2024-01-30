package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	max     = 10
	endTime = 50 * time.Second
)

type Shop struct {
	place   chan bool
	endTime time.Time
	wg      sync.WaitGroup
}

func newshop() *Shop {
	return &Shop{
		place:   make(chan bool, max),
		endTime: time.Now().Add(endTime),
	}
}

func (bs *Shop) Barber() {
	for {
		select {
		case <-time.After(1 * time.Second):
			if len(bs.place) > 0 {
				<-bs.place
				fmt.Println("Barber is working.")
			} else if time.Now().After(bs.endTime) {
				fmt.Println("Going home.")
				bs.wg.Done()
				return
			} else {
				fmt.Println("No clients")
			}
		}
	}
}

func (bs *Shop) Client() {
	defer bs.wg.Done()
	select {
	case bs.place <- true:
		fmt.Println("Client enters.")
	default:
		fmt.Println("place is full. Client leaves.")
	}
}

func main() {
	bs := newshop()

	bs.wg.Add(1)
	go bs.Barber()

	startTime := time.Now()
	for time.Now().Before(startTime.Add(endTime)) {
		time.Sleep(2 * time.Second)
		bs.wg.Add(1)
		go bs.Client()
	}

	bs.wg.Wait()
	fmt.Println("All clients served.")
}
