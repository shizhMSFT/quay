package quay_test

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/shizhMSFT/quay"
	"golang.org/x/exp/constraints"
)

type Slice[T constraints.Ordered] []any

func (s Slice[T]) Len() int           { return len(s) }
func (s Slice[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Slice[T]) Less(i, j int) bool { return s[i].(T) < s[j].(T) }

func ExampleWharf() {
	wharf := quay.NewWharf()
	var wg sync.WaitGroup

	fmt.Println("Checkpoint 1")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if captain := <-wharf.Enter(i); captain {
				time.Sleep(100 * time.Millisecond)
				tickets := wharf.Close()
				time.Sleep(100 * time.Millisecond)
				sort.Sort(Slice[int](tickets))
				fmt.Println(tickets)
				wharf.Arrive()
			}
		}(i)
	}

	time.Sleep(150 * time.Millisecond)

	fmt.Println("Checkpoint 2")
	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if captain := <-wharf.Enter(i); captain {
				time.Sleep(100 * time.Millisecond)
				tickets := wharf.Close()
				sort.Sort(Slice[int](tickets))
				fmt.Println(tickets)
				wharf.Arrive()
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

func ExampleQuay() {
	quay := quay.New()
	var wg sync.WaitGroup

	fmt.Println("Checkpoint 1")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if wharf, captain, dispose := quay.Enter("#1", i); <-captain {
				defer dispose()
				time.Sleep(100 * time.Millisecond)
				tickets := wharf.Close()
				sort.Sort(Slice[int](tickets))
				fmt.Println("#1", tickets)
				wharf.Arrive()
			}
		}(i)
	}

	fmt.Println("Checkpoint 2")
	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if wharf, captain, dispose := quay.Enter("#2", i); <-captain {
				defer dispose()
				time.Sleep(50 * time.Millisecond)
				tickets := wharf.Close()
				sort.Sort(Slice[int](tickets))
				fmt.Println("#2", tickets)
				wharf.Arrive()
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
