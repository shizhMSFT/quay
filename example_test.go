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
			wharf.Travel(i, func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}, func(tickets []int) error {
				time.Sleep(100 * time.Millisecond)
				slices.Sort(tickets)
				fmt.Println(tickets)
				return nil
			})
		}(i)
	}

	time.Sleep(150 * time.Millisecond)

	fmt.Println("Checkpoint 2")
	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			wharf.Travel(i, func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}, func(tickets []int) error {
				slices.Sort(tickets)
				fmt.Println(tickets)
				return nil
			})
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

func ExampleWharf_resign() {
	wharf := quay.NewWharf[int]()
	var wg sync.WaitGroup

	fmt.Println("Checkpoint 1")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			wharf.Travel(i, func() error {
				if i == 0 {
					time.Sleep(50 * time.Millisecond)
					return errors.New("resign")
				}
				time.Sleep(100 * time.Millisecond)
				return nil
			}, func(tickets []int) error {
				time.Sleep(100 * time.Millisecond)
				slices.Sort(tickets)
				fmt.Println(tickets)
				return nil
			})
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
			wharf.Travel(i, func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}, func(tickets []int) error {
				slices.Sort(tickets)
				fmt.Println(tickets)
				return nil
			})
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

func ExampleWharf_arrive() {
	wharf := quay.NewWharf[int]()
	var wg sync.WaitGroup

	fmt.Println("Checkpoint 1")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := wharf.Travel(i, func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}, func(tickets []int) error {
				time.Sleep(100 * time.Millisecond)
				if i == 0 {
					return errors.New("abandon ship")
				}
				slices.Sort(tickets)
				fmt.Println(tickets)
				return nil
			})
			fmt.Println(err)
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
			wharf.Travel(i, func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}, func(tickets []int) error {
				slices.Sort(tickets)
				fmt.Println(tickets)
				return nil
			})
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
			quay.Travel("#1", i, func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}, func(tickets []int) error {
				slices.Sort(tickets)
				fmt.Println("#1", tickets)
				return nil
			})
		}(i)
	}

	fmt.Println("Checkpoint 2")
	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			quay.Travel("#2", i, func() error {
				time.Sleep(50 * time.Millisecond)
				return nil
			}, func(tickets []int) error {
				slices.Sort(tickets)
				fmt.Println("#2", tickets)
				return nil
			})
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
