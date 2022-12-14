package quay

import "sync"

// Quay is a quay with scalable wharves.
type Quay[T any] struct {
	gate    sync.Mutex
	wharves map[any]*Wharf[T]
}

// New creates a virtual scalable quay.
func New[T any]() *Quay[T] {
	return &Quay[T]{}
}

// Travel travels to the destination through a specific wharf, holding a ticket.
func (q *Quay[T]) Travel(wharfID any, ticket T, learn func() error, sail func(tickets []T) error) error {
	wharf := q.wharf(wharfID)
	defer func() {
		q.gate.Lock()
		defer q.gate.Unlock()
		if wharf.idle() {
			delete(q.wharves, wharfID)
		}
	}()
	return wharf.Travel(ticket, learn, sail)
}

func (q *Quay[T]) wharf(wharfID any) *Wharf[T] {
	q.gate.Lock()
	defer q.gate.Unlock()

	wharf, ok := q.wharves[wharfID]
	if !ok {
		wharf = NewWharf[T]()
		if q.wharves == nil {
			q.wharves = map[any]*Wharf[T]{
				wharfID: wharf,
			}
		} else {
			q.wharves[wharfID] = wharf
		}
	}
	return wharf
}
