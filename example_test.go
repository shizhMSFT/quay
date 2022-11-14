package quay_test

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/shizhMSFT/quay"
	"golang.org/x/exp/slices"
)

func ExampleWharf() {
	wharf := quay.NewWharf[int]()
	var wg sync.WaitGroup

	fmt.Println("Checkpoint 1")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if status := <-wharf.Enter(i); status.Elected {
				time.Sleep(100 * time.Millisecond)
				tickets := wharf.Close()
				time.Sleep(100 * time.Millisecond)
				slices.Sort(tickets)
				fmt.Println(tickets)
				wharf.Arrive(nil)
			}
		}(i)
	}

	time.Sleep(150 * time.Millisecond)

	fmt.Println("Checkpoint 2")
	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if status := <-wharf.Enter(i); status.Elected {
				time.Sleep(100 * time.Millisecond)
				tickets := wharf.Close()
				slices.Sort(tickets)
				fmt.Println(tickets)
				wharf.Arrive(nil)
			}
		}(i)
	}

	fmt.Println("Checkpoint 3")
	wg.Wait()

	// Output:
	// Checkpoint 1
	// Checkpoint 2
	// Checkpoint 3
	// [0 1 2 3 4]
	// [5 6 7 8 9]
}

func ExampleWharf_Resign() {
	wharf := quay.NewWharf[int]()
	var wg sync.WaitGroup

	fmt.Println("Checkpoint 1")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if status := <-wharf.Enter(i); status.Elected {
				if i == 0 {
					time.Sleep(50 * time.Millisecond)
					wharf.Resign()
					return
				}
				time.Sleep(100 * time.Millisecond)
				tickets := wharf.Close()
				time.Sleep(100 * time.Millisecond)
				slices.Sort(tickets)
				fmt.Println(tickets)
				wharf.Arrive(nil)
			}
		}(i)
		if i == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	time.Sleep(150 * time.Millisecond)

	fmt.Println("Checkpoint 2")
	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if status := <-wharf.Enter(i); status.Elected {
				time.Sleep(100 * time.Millisecond)
				tickets := wharf.Close()
				slices.Sort(tickets)
				fmt.Println(tickets)
				wharf.Arrive(nil)
			}
		}(i)
	}

	fmt.Println("Checkpoint 3")
	wg.Wait()

	// Output:
	// Checkpoint 1
	// Checkpoint 2
	// Checkpoint 3
	// [1 2 3 4]
	// [5 6 7 8 9]
}

func ExampleWharf_Arrive() {
	wharf := quay.NewWharf[int]()
	var wg sync.WaitGroup

	fmt.Println("Checkpoint 1")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			status := <-wharf.Enter(i)
			if status.Elected {
				time.Sleep(100 * time.Millisecond)
				tickets := wharf.Close()
				time.Sleep(100 * time.Millisecond)
				if i == 0 {
					wharf.Arrive(errors.New("abandon ship"))
					return
				}
				slices.Sort(tickets)
				fmt.Println(tickets)
				wharf.Arrive(nil)
			} else {
				fmt.Println(status.Error)
			}
		}(i)
		if i == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	time.Sleep(150 * time.Millisecond)

	fmt.Println("Checkpoint 2")
	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if status := <-wharf.Enter(i); status.Elected {
				time.Sleep(100 * time.Millisecond)
				tickets := wharf.Close()
				slices.Sort(tickets)
				fmt.Println(tickets)
				wharf.Arrive(nil)
			}
		}(i)
	}

	fmt.Println("Checkpoint 3")
	wg.Wait()

	// Output:
	// Checkpoint 1
	// Checkpoint 2
	// Checkpoint 3
	// abandon ship
	// abandon ship
	// abandon ship
	// abandon ship
	// [5 6 7 8 9]
}

func ExampleQuay() {
	quay := quay.New[int]()
	var wg sync.WaitGroup

	fmt.Println("Checkpoint 1")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if wharf, status, dispose := quay.Enter("#1", i); (<-status).Elected {
				defer dispose()
				time.Sleep(100 * time.Millisecond)
				tickets := wharf.Close()
				slices.Sort(tickets)
				fmt.Println("#1", tickets)
				wharf.Arrive(nil)
			}
		}(i)
	}

	fmt.Println("Checkpoint 2")
	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if wharf, status, dispose := quay.Enter("#2", i); (<-status).Elected {
				defer dispose()
				time.Sleep(50 * time.Millisecond)
				tickets := wharf.Close()
				slices.Sort(tickets)
				fmt.Println("#2", tickets)
				wharf.Arrive(nil)
			}
		}(i)
	}

	fmt.Println("Checkpoint 3")
	wg.Wait()

	// Output:
	// Checkpoint 1
	// Checkpoint 2
	// Checkpoint 3
	// #2 [5 6 7 8 9]
	// #1 [0 1 2 3 4]
}
